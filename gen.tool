version: "0.1"
database:
  # consult[https://gorm.io/docs/connecting_to_the_database.html]"
  dsn : "root:root@tcp(127.0.0.1:3310)/min_mall?charset=utf8mb4&parseTime=true&loc=Local"
  # input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
  db  : "mysql"
  # specify a directory for output
  outPath :  "./dao/query"
  # query code file name, default: gen.go
  outFile :  ".go"
  # generate unit test for query code
  withUnitTest  : false
  # generated model code's package name
  modelPkgName  : ""
  # generate with pointer when field is nullable
  fieldNullable : false
  # generate field with gorm index tag
  fieldWithIndexTag : false
  # generate field with gorm column type tag
  fieldWithTypeTag  : false
  # only generate model with out query
  onlyModel : true