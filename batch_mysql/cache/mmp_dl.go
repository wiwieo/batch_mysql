// +build linux,cgo darwin,cgo

package cache

import (
	"os"
	"syscall"
	"time"
	"wiwieo/batch_mysql/constant"
)

type mmap struct {
	data         []byte        // 与文件映射的内存
	dataC        chan []byte   // 用于写入的通道
	stopC        chan struct{} // 停止写入
	stop         bool
	f            *os.File // 日志文件
	FullFilePath string   // 文件全路径
	FilePath     string   // 文件路径
	FileName     string   // 文件名称
	at           int      // 在什么位置写
	size         int      // 与文件映射的大小
	initSize     int      // 与文件映射的初始大小
	isM          bool     // 如果使用mmap出错，则改用直接写文件的方式
	w            bool     // 正在清理，等待写入
}

func NewMmap(filePath, name string, size int) (*mmap, error) {
	// 构建对应的结构体，以配后续使用
	m := &mmap{
		size:         size,
		FullFilePath: filePath,
		FilePath:     filePath,
		dataC:        make(chan []byte, 10),
		stopC:        make(chan struct{}, 1),
		stop:         false,
		isM:          true,
		initSize:     size,
		FileName:     name,
	}

	// 使用channel方式，同步写入
	go m.wait()
	//go m.rename()
	return m, m.init(filePath)
}

func (m *mmap) Read() (content []byte, err error) {
	// open a file
	// get file size
	fsize, err := syscall.Seek(int(m.f.Fd()), 0, 2)
	if err != nil {
		return nil, err
	}
	content = make([]byte, fsize)
	content, err = syscall.Mmap(int(m.f.Fd()), 0, int(fsize), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return
}

// 只有在文件中有实质内容时，才进行分隔
// 修改日志文件名称，用于分割日志使用，防止一个文件过大
func (m *mmap) rename(size int) {
	if m.at > 0 {
		m.unmap()
		err := os.Rename(m.FullFilePath, m.FullFilePath+"."+time.Now().Format("20060102150405"))
		if err == nil {
			m.size = size
			m.allocate()
		}
	}
}

// 初始化log信息
func (m *mmap) init(filePath string) error {
	filePath += string(os.PathSeparator) + m.FileName
	os.MkdirAll(filePath, os.ModePerm)
	filePath += string(os.PathSeparator) + m.FileName + constant.Suffix
	m.FullFilePath = filePath
	// 存在，则改名
	os.Rename(filePath, filePath+"."+time.Now().Format("20060102150405"))
	err := m.setFileInfo(filePath)
	if err != nil {
		return err
	}

	err = m.allocate()
	if err != nil {
		return err
	}
	return nil
}

// MMAP映射
func (m *mmap) allocate() error {
	if m.f == nil {
		m.setFileInfo(m.FullFilePath)
	}
	defer func() {
		m.f.Close()
		m.f = nil
	}()

	// 文件映射的大小必须是页数的倍数，如果不是，则自动根据大小调整为相应倍数
	if m.size%syscall.Getpagesize() != 0 {
		m.size = (m.size / syscall.Getpagesize()) * syscall.Getpagesize()
	}
	if m.size == 0 {
		m.size = syscall.Getpagesize()
	}

	// MMAP映射时，文件必须有相应大小的内容，即需要相应大小的占位符
	if _, err := m.f.WriteAt(make([]byte, m.size), int64(m.at)); nil != err {
		return err
	}

	// 映射
	data, err := syscall.Mmap(int(m.f.Fd()), 0, int(m.size), syscall.PROT_WRITE|syscall.PROT_READ|syscall.PROT_EXEC, syscall.MAP_SHARED)
	if nil != err {
		return err
	}
	m.data = data
	return nil
}

// 设置映射的文件
func (m *mmap) setFileInfo(filePath string) error {
	// 打开文件，不存在创建新文件
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if nil != err {
		return err
	}

	// 获取当前文件信息
	fi, err := f.Stat()
	if nil != err {
		return err
	}

	m.f = f
	m.at = int(fi.Size())
	return nil
}

// 关闭所有
func (m *mmap) Close() error {
	m.stopC <- struct{}{}
	// 需要时间去处理后续操作，包括未写入的数据
	time.Sleep(100 * time.Millisecond)
	return nil
}

// 关闭文件映射
func (m *mmap) unmap() error {
	// 关闭映射
	if err := syscall.Munmap(m.data); nil != err {
		return err
	}

	// 将未写入的内容清空
	// 如果未清空，在文件末位未写入位置，将会出现大量占位符
	err := os.Truncate(m.FullFilePath, int64(m.at))
	if err != nil {
		return err
	}
	return nil
}

// 将已经写入的内容清除
func (m *mmap) Clean() error {
	m.w = true
	defer func() {
		m.at = 0
		m.w = false
	}()
	for i := 0; i < m.at; i++ {
		m.data[i] = ' '
	}
	return nil
}

// 接收写入内容
func (m *mmap) Write(content []byte) error {
	if !m.stop {
		m.dataC <- content
	}
	return nil
}

// 当初始映射大小不足时，进行扩容
func (m *mmap) doubleAllocate() error {
	// 先将之前的映射关闭
	m.unmap()
	m.size = 2 * m.size
	return m.allocate()
}

// 等待内容写入
func (m *mmap) wait() {
	//t := time.NewTimer(time.Duration(GetTimeer(time.Now())))
	for {
		select {
		case content, ok := <-m.dataC:
			// 通道被关闭且服务停止，则关闭映射
			if !ok && m.stop {
				m.unmap()
				return
			}
			if len(content) == 0 {
				return
			}
			// 等待清理完成
			for m.w {
			}
			// 剩余空间不足以添加所有内容，需要扩容
			for len(content) > m.size-m.at {
				err := m.doubleAllocate()
				if err != nil {
					m.isM = false
					m.unmap()
				}
			}
			m.write(content)
		//case <-t.C:
		//	// 每天更新一次名称，用于分隔文件
		//	// 考虑到每天的量应该差不多，故此处新文件大小，直接在原文件的1/2
		//	m.rename(m.at / 2)
		//	t.Reset(time.Duration(GetTimeer(time.Now())))

		case <-m.stopC:
			// 停止往channel里继续写数据
			m.stop = true
			// 关闭channel
			close(m.dataC)
		}
	}
}

func (m *mmap) write(content []byte) {
	if m.isM {
		m.writeWithMmap(content)
	} else {
		m.writeWithIO(content)
	}
}

func (m *mmap) writeWithMmap(content []byte) {
	// 内容写入文件
	for i, v := range content {
		m.data[m.at+i] = v
	}
	m.at += len(content)
}

func (m *mmap) writeWithIO(content []byte) {
	var err error
	if m.f == nil {
		m.f, err = os.OpenFile(m.FullFilePath, os.O_RDWR|os.O_CREATE, os.ModeAppend)
		if err != nil {
			panic(err)
		}
	}

	size, err := m.f.WriteAt(content, int64(m.at))
	m.at += size
	if err != nil {
		panic(err)
	}
}

func GetTimeer(now time.Time) int64 {
	dest := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 5, 0, 0, time.Local)
	return dest.UnixNano() - now.UnixNano()
}
