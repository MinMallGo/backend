package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func mkdir() {
	s := []string{
		"api",
		"cert",
		"config/config.go",
		"constant/constants.go",
		"custom_struct/custom.go",
		"dao/user_dao.go",
		"database/database.go",
		"dto/user_dto.go",
		"entity/user.go",
		"enumerate/status_enum.go",
		"middleware/auth_middleware.go",
		"response/response.go",
		"route/routes.go",
		"service/user_service.go",
		"test/user_test.go",
		"util/hash_util.go",
		"vo/user_vo.go",
		"main.go",
	}

	for _, s2 := range s {
		val := strings.Split(s2, "/")
		err := os.Mkdir(val[0], 0644)
		if err != nil {
			log.Fatal("mkdir: ", val[0], " with error:", err)
		}

		if len(val) > 1 {
			val2 := val[1:]
			filepath := val[0]
			for i, v := range val2 {
				if i == len(val2)-1 {
					// 创建文件
					filepath = path.Join(filepath, v)
					create, err := os.Create(filepath)
					if err != nil {
						log.Fatal("mkfile: ", filepath, " with error:", err)
					}
					v = strings.TrimSuffix(v, ".go")
					create.Write([]byte("package " + v))
					create.Close()
				} else {
					// 创建文件夹
					filepath = path.Join(filepath, v)
					err := os.Mkdir(path.Dir(filepath), 0644)
					if err != nil {
						log.Fatal("mkdir: ", filepath, " with error:", err)
					}
				}
			}
		}
	}
}
