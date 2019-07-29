package cache

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"sync"
	"time"
	"wiwieo/batch_mysql/constant"
)

type Cache struct {
	sync.RWMutex
	content   [][]byte      // 内容
	totalSize uint32        // 总容量，目前没有用到
	lastTime  time.Time     // 上次持久时间
	persist   chan struct{} // 持久
	watchIt   chan struct{} // 检查是否触发持久条件
	isBackup  bool          // 批量写入数据库之前，是否将缓存的数据先持久到硬盘中，以防止服务崩溃时数据不丢失
	m         Backup        // 如果isBackup是true，则使用mmap方式写入文件；false，则无效
	w         Write
	name      string
}

type Write interface {
	WriteToMySQL([][]byte) error
}

// w: 实现Write接口的对象，在触发时调用，用于向数据库写入数据
func NewCache(w Write, name string) (*Cache, error) {
	c := &Cache{
		totalSize: 100,
		lastTime:  time.Now(),
		persist:   make(chan struct{}, 2),
		watchIt:   make(chan struct{}, 1),
		w:         w,
		isBackup:  constant.Config.IsPersist,
		content:   make([][]byte, 0, constant.Config.MaxTriggerSize*2),
		name:      name,
	}
	if constant.Config.IsPersist {
		m, err := NewMmap(constant.Config.Path, name, 1<<10)
		c.m = m
		if err != nil {
			c.m = nil
		}
	}
	go c.watch()
	go c.write()
	go c.copeHistoryData()
	return c, nil
}

func (c *Cache) watch() {
	t := time.Tick(time.Duration(constant.Config.MaxTriggerTime+2) * time.Second)
	for {
		select {
		case <-c.watchIt:
			c.watchSize()
		case <-t:
			c.watchTime()
		}
	}
}

func (c *Cache) watchTime() {
	if time.Since(c.lastTime) >= time.Duration(constant.Config.MaxTriggerTime)*time.Second {
		c.persist <- struct{}{}
	}
}

func (c *Cache) watchSize() {
	if len(c.content) >= constant.Config.MaxTriggerSize {
		c.persist <- struct{}{}
	}
}

func (c *Cache) AddContent(content []byte) error {
	c.Lock()

	if c.m != nil {
		var temp = make([]byte, len(content))
		copy(temp, content)
		temp = append(temp, constant.Separate...)
		err := c.m.Write(temp)
		if err != nil {
			c.Unlock()
			return err
		}
	}
	c.content = append(c.content, content)
	c.Unlock()
	// 不能放在锁内，容易造成死锁
	c.watchIt <- struct{}{}
	return nil
}

func (c *Cache) write() {
	for {
		select {
		case <-c.persist:
			if len(c.content) == 0 {
				continue
			}
			// TODO 是否可以不使用锁来处理
			c.Lock()
			go func(content [][]byte) {
				err := c.w.WriteToMySQL(content)
				if err != nil {
					glog.Errorf("write to mysql is wrong, %s\nthe content filed is: %s", err, content)
				}
			}(c.content)
			c.content = make([][]byte, 0)
			if c.m != nil {
				err := c.m.Clean()
				if err != nil {
					c.Unlock()
					glog.Errorf("clean backup file is wrong, %s", err)
					continue
				}
			}
			c.Unlock()
			c.lastTime = time.Now()
		}
	}
}

func (c *Cache) copeHistoryData() {
	contents := c.readDataFromPath(constant.Config.Path)
	for _, cs := range contents {
		for _, content := range cs {
			c.AddContent(content)
		}
	}
}

// [文件][行数][内容]
func (c *Cache) readDataFromPath(path string) [][][]byte {
	path += string(os.PathSeparator) + c.name + string(os.PathSeparator)
	// 文件夹不存在，则说明没有需要读取的文件，直接跳过
	fs, _ := ioutil.ReadDir(path)
	rtn := make([][][]byte, 0, len(fs))
	for _, f := range fs {
		// 跳过目录
		if f.IsDir() {
			continue
		}
		if f.Name() == c.name+constant.Suffix {
			continue
		}
		fileName := path + string(os.PathSeparator) + f.Name()
		rtn = append(rtn, c.readDataFromFile(fileName))
		os.Remove(fileName)
	}
	return rtn
}

func (c *Cache) readDataFromFile(fileName string) [][]byte {
	contents, _ := ioutil.ReadFile(fileName)
	var rtn = make([][]byte, 0, constant.Config.MaxTriggerSize)
	idxs := c.findSeparateIdx(contents)
	start := 0
	for _, idx := range idxs {
		rtn = append(rtn, contents[start:idx])
		start = idx + len(constant.Separate)
	}
	return rtn
}

func (c *Cache) findSeparateIdx(content []byte) []int {
	idxs := make([]int, 0, constant.Config.MaxTriggerSize)
	for idx := 0; idx < len(content)-len(constant.Separate)+1; idx++ {
		for i, s := range constant.Separate {
			if content[idx+i] != s {
				break
			}
			if i == len(constant.Separate)-1 {
				idxs = append(idxs, idx)
			}
		}
	}
	return idxs
}
