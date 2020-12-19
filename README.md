# Gorm库

- Full-Featured ORM
- Associations (Has One, Has Many, Belongs To, Many To Many, Polymorphism, Single-table inheritance)
- Hooks (Before/After Create/Save/Update/Delete/Find)
- Eager loading with `Preload`, `Joins`
- Transactions, Nested Transactions, Save Point, RollbackTo to Saved Point
- Context, Prepared Statment Mode, DryRun Mode
- Batch Insert, FindInBatches, Find/Create with Map, CRUD with SQL Expr and Context Valuer
- SQL Builder, Upsert, Locking, Optimizer/Index/Comment Hints, Named Argument, SubQuery
- Composite Primary Key, Indexes, Constraints
- Auto Migrations
- Logger
- Extendable, flexible plugin API: Database Resolver (Multiple Databases, Read/Write Splitting) / Prometheus…
- Every feature comes with tests
- Developer Friendly

## 入门

###　安装

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite   # for sqlite
go get -u gorm.io/driver/mysql    # for MySQL
```

### 快速入门

```go
package main

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

type Product struct {
  gorm.Model
  Code  string
  Price uint
}

func main() {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  db.AutoMigrate(&Product{})

  // Create
  db.Create(&Product{Code: "D42", Price: 100})

  // Read
  var product Product
  db.First(&product, 1) // find product with integer primary key
  db.First(&product, "code = ?", "D42") // find product with code D42

  // Update - update product's price to 200
  db.Model(&product).Update("Price", 200)
  // Update - update multiple fields
  db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
  db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

  // Delete - delete product
  db.Delete(&product, 1)
}
```

> **注意**：
>
> 	1. `db.Update` 和 `db.Updates` 两个函数不同，前者**只能更新一列**的属性值
> 	2. `db.Updates` 函数使用结构体作为参数时，**会忽略为0，NULL的字段**
>
> 2. 如果采用内嵌 `gorm.Model`的结构体, `db.Delete`函数**并不会从数据中将其删除**，而是将表中`delete_at`字段由原来的空值设置成当前的时间戳

```sql
mysql> SELECT * FROM products;
+----+-------------------------+-------------------------+-------------------------+------+-------+
| id | created_at              | updated_at              | deleted_at              | code | price |
+----+-------------------------+-------------------------+-------------------------+------+-------+
|  1 | 2020-12-18 21:36:06.058 | 2020-12-18 21:36:06.286 | 2020-12-18 21:36:06.327 | F42  |   200 |
|  2 | 2020-12-18 21:39:58.085 | 2020-12-18 21:39:58.234 | NULL                    | F42  |   200 |
|  3 | 2020-12-18 21:42:51.107 | 2020-12-18 21:46:34.550 | 2020-12-18 21:46:34.594 | E42  |   200 |
|  4 | 2020-12-18 21:44:23.888 | 2020-12-18 21:44:23.888 | NULL                    | E42  |   100 |
|  5 | 2020-12-18 21:45:03.390 | 2020-12-18 21:45:03.390 | NULL                    | E42  |   100 |
|  6 | 2020-12-18 21:46:34.495 | 2020-12-18 21:46:34.495 | NULL                    | E42  |   100 |
+----+-------------------------+-------------------------+-------------------------+------+-------+
6 rows in set (0.00 sec)
```





### 定义模型

建模就是采用`GO`语言中的结构体

```go
type User struct {
  ID           uint
  Name         string
  Email        *string
  Age          uint8
  Birthday     *time.Time
  MemberNumber sql.NullString
  ActivedAt    sql.NullTime
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
```



### 转换

`GORM`首选配置，默认`ID` 最为主键, 将结构的名称复数化 `snake_cases` 做为表名称, 而`snake_case` 作为列名称 ，并且使用 `CreatedAt`, `UpdatedAt` 跟踪创建、更新的时间。

> 注意：表和列的名称采用驼峰式将转换成`word1_word2`形式,无论你定义的结构体以及结构体中的字段是否大写，最终在数据库中呈现的都为小写的名称，如果采用匈牙利命名法不仅**小写所有字母**并按照大写字母进行分割，**用`_`下划线衔接**。



### gorm.Model

你可以将`gorm.Model`嵌入到其它模型中

```go
// gorm.Model definition
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```



### 字段级别的权限

默认到处的字段具有`CRUD`的全部权限，比较有意思的是`GORM`可以让你控制读取和改写的权限

`->`指向右侧的箭头代表`read`

`<-`指向左侧的箭头则代表`write`和`update`

箭头指向之后可以附带修改的权限，如果为`:false`代表禁用此权限

* :create 
* :update 
* :false

```go
type User struct {
  Name string `gorm:"<-:create"` // allow read and create
  Name string `gorm:"<-:update"` // allow read and update
  Name string `gorm:"<-"`        // allow read and write (create and update)
  Name string `gorm:"<-:false"`  // allow read, disable write permission
  Name string `gorm:"->"`        // readonly (disable write permission unless it configured )
  Name string `gorm:"->;<-:create"` // allow read and create
  Name string `gorm:"->:false;<-:create"` // createonly (disabled read from db)
  Name string `gorm:"-"`  // ignore this field when write and read
}
```



### 时间格式和精度

```go
type User struct {
  CreatedAt time.Time // Set to current time if it is zero on creating
  UpdatedAt int       // Set to current unix seconds on updaing or if it is zero on creating
  Updated   int64 `gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
  Updated   int64 `gorm:"autoUpdateTime:milli"`// Use unix milli seconds as updating time
  Created   int64 `gorm:"autoCreateTime"`      // Use unix seconds as creating time
}
```



### 内嵌结构体

```go
type User struct {
  gorm.Model
  Name string
}
// equals
type User struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
  Name string
}
```

对于普通的结构体,可以添加注释体`gorm:"embedded"`

```go
type Author struct {
  Name  string
  Email string
}

type Blog struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
// equals
type Blog struct {
  ID    int64
  Name  string
  Email string
  Upvotes  int32
}
```

还可以使用标签`embeddedPrefix`添加字段的名字的前缀

```go
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded;embeddedPrefix:author_"`
  Upvotes int32
}
// equals
type Blog struct {
  ID          int64
  AuthorName  string
  AuthorEmail string
  Upvotes     int32
}
```



> **注意：** 同样定义结构体，如果你使用了内嵌`gorm.Model`,在执行删除操作时不会真正的将数据库中的数据删除，而只是将`delete_at`设置为操作的时间戳。同样在查询中，默认会自动添加`WHERE delete_at IS NOT NULL`，因而让我看不到查询标记为删除的数据。



### 字段标签

[其他字段标签属性](https://gorm.io/docs/models.html#Fields-Tags)



### 连接数据库

MySQL

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

```

> **注意:**
> 如果要处理 `time.Time` , 你需要包含 `parseTime` 最为参数. ([更多参数](https://github.com/go-sql-driver/mysql#parameters))
> 要完全支持 UTF-8 编码，你需要将 `charset=utf8` 改为 `charset=utf8mb4`. 




-------



## 查询

### 创建

#### 添加条目

```go
now := time.now()
user := User{Name: "Jinzhu", Age: 18, Birthday: &now}

result := db.Create(&user) // pass pointer of data to Create

user.ID             // returns inserted data's primary key
result.Error        // returns error
result.RowsAffected // returns inserted records count
```

> **注意：** GORM所使用的为 \*time.Time,如果要取当前时候需要将`time.Now`赋予给变量，再取变量的地址



#### 使用所选择的字段创建条目

1. 根据制定的字段添加一条记录，先使用`db.Select`选择所需要操作的列，然后再使用`db.Create`函数

```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775"
```

> 纵使你初始化的结构体中有着每个属性和其对应值，但是使用 `db.Select`只会将特定的列插入到数据库

2. 如果表的字段过多，也可以用排除法，去除某些字段插入记录, 使用`db.Omit`函数

```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")
```

3. `db.Create`也可以对slice数据批量的插入

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

for _, user := range users {
    fmt.Println(user.ID) // 1,2,3
}
```

4. 另外也可以使用`db.CreateInBatches`第二个参数**指定批量插入的大小**创建记录

```go
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)
```



#### 数据创建时的钩子

1. GORM允许用户定义钩子功能，`BeforeSave`, `BeforeCreate`, `AfterSave`, `AfterCreate`。这些钩子函数对应数据创建的不同阶段

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New()

  if u.Role == "admin" {
    return errors.New("invalid role")
  }
  return
}
```

> 注意：如果对结构体做任何的修改，你需要使用`db.AutoMigrate()`将所做的变化反映到数据库中

2. 如果你想跳过钩子函数，你可以`grom.Session`使用`SkipHooks`属性

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)

DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)

```

#### 以Map创建条目

GORM 支持从`map[string]interface{}`和`[]map[string]interface{}`中创建条目。

```go
db.Model(&User{}).Create(map[string]interface{}{
  "Name": "jinzhu", "Age": 18,
})

// batch insert from `[]map[string]interface{}{}`
db.Model(&User{}).Create([]map[string]interface{}{
  {"Name": "jinzhu_1", "Age": 18},
  {"Name": "jinzhu_2", "Age": 20},
})
```



#### 使用SQL函数表达式

在某些应用场景，你可以需要调用到数据库内部集成的函数或者存储过程表达式，GORM通过`clause.Expr`结构体让我们也可以自由的形成SQL语句

```go
// This is example intend to call `UPPER` function which builtin MySQL

package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User Model
type SqlExpressUser struct {
	Name     string
	Location string
}

func main() {
	dsn := "root:redhat@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to MySQL server.")
	}

	db.AutoMigrate(&SqlExpressUser{})

	db.Model(&SqlExpressUser{}).Create(map[string]interface{}{
		"Name": "JinZhu",
		"Location": clause.Expr{
			SQL:  "UPPER(?)",
			Vars: []interface{}{"HuNan ChangSha"},
		},
	})
}

```

[更为复杂的用法请参考](https://gorm.io/docs/create.html#Create-From-SQL-Expression-Context-Valuer)

> 几个要注意的地方，`clause.Expr{}`结构体中`SQL`多为SQL函数，因此别忘了函数的调用括号和传参`(?)`。
>
> 第二 `Vars`数据类型为`[]inteface{}`,这个切片你可以根据问号的多少对应传具体的值



#### 关联型的数据添加

当创建一些相关联的数据，如果相关联的数据非零和空值，这些关联的数据同样保持更新，同样它的钩子一样会被执行

```go
type AssoCreditCard struct {
  gorm.Model
  Number   string
  AssoUserID   uint
}

type AssoUser struct {
  gorm.Model
  Name       string
  CreditCard AssoCreditCard
}

db.Create(&AssoUser{
  Name: "jinzhu",
  CreditCard: AssoCreditCard{Number: "411111111111"}
})
// INSERT INTO `asso_users` ...
// INSERT INTO `asso_credit_cards` ...
```

> 注意上面的例子中 `AssoCreditCard`中定义的`AssoUserID`,GORM确定关联关系有着自己的规则，上面的例子使用`Belong to`一对一的关系，`AssoUserID`分成两部分，`AssoUser`为关联的结构体，`ID`一定要全大写，这样就确定了关系，更多请查考下面的内容，另外请参考[官方文档](https://gorm.io/docs/belongs_to.html#Belongs-To)



#### 默认值

```go
type User struct {
  ID   int64
  Name string `gorm:"default:galeone"`
  Age  int64  `gorm:"default:18"`
}
```

```go
type User struct {
  gorm.Model
  Name string
  Age  *int           `gorm:"default:18"`
  Active sql.NullBool `gorm:"default:true"`
}
```

```go
type User struct {
  ID        string `gorm:"default:uuid_generate_v3()"` // db func
  FirstName string
  LastName  string
  Age       uint8
  FullName  string `gorm:"->;type:GENERATED ALWAYS AS (concat(firstname,' ',lastname));default:(-);`
}
```



#### 数据冲突

GORM 对多种数据实现的冲突数据的处理

```go
import "gorm.io/gorm/clause"

// Do nothing on conflict
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// Update columns to default value on `id` conflict
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE ***; MySQL

// Update columns to new value on `id` conflict
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age=VALUES(age); MySQL

// Update all columns expects primary keys to new value on conflict
db.Clauses(clause.OnConflict{
  UpdateAll: true
}).Create(&users)
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age", ...;
```





### Struct 和 Map 条件查询

```go
// 使用结构体
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// 使用 Map
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// 使用切片的主键
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

> GORM 仅仅查询非零的字段,下面的例子可以看到忽略`Age`为零的条件查询：

```go
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu";
```

> 你可以使用 map 编译查询解决上面的问题

```go
db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
```


### 内置条件

与 `Where` 相同

```go
// SELECT * FROM users WHERE id = 23;
// Get by primary key if it were a non-integer type
db.First(&user, "id = ?", "string_primary_key")
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";

db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

// Struct
db.Find(&users, User{Age: 20})
// SELECT * FROM users WHERE age = 20;

// Map
db.Find(&users, map[string]interface{}{"age": 20})
// SELECT * FROM users WHERE age = 20;
```



#### Not 条件

 使用 Not() 函数实现逻辑取反

```go
db.Not("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

// Not In
db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

// Struct
db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

// Not In slice of primary keys
db.Not([]int64{1,2,3}).First(&user)
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
```



#### Or 条件

```go
b.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

// Map
db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
```



### 裸SQL 和 SQL 编译器



## 关系型数据库

### Belong to 隶属于

1. 隶属于关系

```go
   // `User` belongs to `Company`, `CompanyID` is the foreign key
type User struct {
     gorm.Model
     Name      string
     CompanyID int
     Company   Company
}
   
type Company struct {
     ID   int
     Name string
}
```

   

2. **覆盖外键**

   可以通过注解改写外键的定义

```go
type User struct {
     gorm.Model
     Name         string
     CompanyRefer int
     Company      Company `gorm:"foreignKey:CompanyRefer"`
     // use CompanyRefer as foreign key
}
   
type Company struct {
     ID   int
     Name string
}
```



3. **覆盖外键参考的列**

   默认GORM定义参考的都为ID列。你可以使用注解定义`references:`指向其他的字段，下面的例子，用户所对应的公司参考的外键为公司的`Code`列，而不是缺省的`ID`。

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID string
  Company   Company `gorm:"references:Code"` // use Code as references
}

type Company struct {
  ID   int
  Code string
  Name string
}
```

​	4. **外键约束**

你可以定义外键的约束，比如删除、更新主表时的联动操作

```go
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
  ID   int
  Name string
}
```



### Has One 一对一

`Has one`特性可以说就是`Belong to` 的一对一，因此和`Belong to`大体相同

```go
// User has one CreditCard, CreditCardID is the foreign key
type User struct {
  gorm.Model
  CreditCard CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```



### Has Many 一对多

`Has many` 关系与`Has one` 不同在于，它可以拥有零或者多个其他模型的结构，所以我们只需要用``slice`表现即可

比如下面的例子中一个用户拥有多张信用卡

```go
// User has many CreditCards, UserID is the foreign key
type User struct {
  gorm.Model
  CreditCards []CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```



### Many To Many 多对多

1. 多对多关系涉及到两个模型，因而需要用另外一种注解方式说明`many2many`，**而注解中`user_languages`为生成的中间表，由它存储和反映多对多的关系**。

```go
// User has and belongs to many languages, `user_languages` is the join table
type User struct {
  gorm.Model
  Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
}
```

当使用 `AutoMigrate` 对`User` 表做数据迁移时，GORM自动对`many2many`关系的表迁移



2. 如果需要反向获取参考表数据时，可以在两个结构体中都使用注解`many2many`

```go
// User has and belongs to many languages, use `user_languages` as join table
type User struct {
  gorm.Model
  Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
  Users []*User `gorm:"many2many:user_languages;"`
}
```

> 注意：反向获取参考数据时，`many2many`定义的名称应该相同, 另外注意back-refer使用是**结构体指针切片**



3. 自连接

```go
type User struct {
  gorm.Model
  Friends []*User `gorm:"many2many:user_friends"`
}

// Which creates join table: user_friends
//   foreign key: user_id, reference: users.id
//   foreign key: friend_id, reference: users.id

```



### 关联模式

### 预先加载





## 日志



#### Logger

```go
newLogger := logger.New(
  log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
  logger.Config{
    SlowThreshold: time.Second,   // Slow SQL threshold
    LogLevel:      logger.Silent, // Log level
    Colorful:      false,         // Disable color
  },
)

// Globally mode
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
  Logger: newLogger,
})

// Continuous session mode
tx := db.Session(&Session{Logger: newLogger})
tx.First(&user)
tx.Model(&user).Update("Age", 18)
```



#### 日志级别

GORM定义的日志级别`Silent`, `Error`, `Warn`, `Info`

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Silent),
})
```



#### 调试

单条语句的调试

```go
db.Debug().Where("name = ?", "jinzhu").First(&User{})
```


