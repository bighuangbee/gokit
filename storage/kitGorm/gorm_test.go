package kitGorm
import (
	"fmt"
	kitLog"github.com/bighuangbee/gokit/log"
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestMysql(t *testing.T) {
	logger := kitLog.NewZapLogger(&kitLog.Options{
		Level:       0,
		ServiceName: "",
		Skip:        0,
		Writer:      nil,
	})

	db, err := New(&Options{
		Address:  "192.168.18.66:23306",
		UserName: "root",
		Password: "Hiscene2022",
		DBName:   "hiar_mozi_device",
		Logger:   Logger{L: log.NewHelper(logger)},
	})
	if err != nil{
		fmt.Println(err)
		return
	}

	type user struct {
		Id int
		Username string
		Nickname string
	}

	for i := 0; i < 10; i++ {
		var u user
		db.Raw("select * from confs order by id asc limit 1").Scan(&u)
		fmt.Println(u)
	}
}
