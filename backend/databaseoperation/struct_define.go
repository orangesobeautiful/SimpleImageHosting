package databaseoperation

// Setting datebase setting struct
type Setting struct {
	Name  string `gorm:"column: Name; type:VARCHAR(30) NOT NULL; primary_key;" json:"name"`
	Value string `gorm:"column: Value; type:TEXT NOT NULL;" json:"value"`
}

// User datebase user struct
type User struct {
	ID            int64  `gorm:"column: ID; type: BIGINT UNSIGNED NOT NULL auto_increment; primary_key;" json:"id"`
	LoginName     string `gorm:"column: LoginName; type:VARCHAR(30) NOT NULL; uniqueIndex:idx_loginname;" json:"login_name"`
	ShowName      string `gorm:"column: ShowName; type:VARCHAR(30) NOT NULL;" json:"show_name"`
	Email         string `gorm:"column: Email; type:VARCHAR(256) NOT NULL; uniqueIndex:idx_email;" json:"email"`
	PwdHash       []byte `gorm:"column: PwdHash; type:BINARY(60) NOT NULL;" json:"pwd_hash"`
	Avatar        string `gorm:"column: Avatar; type:VARCHAR(30) NOT NULL; default:\"\"" json:"avatar"`
	Introduction  string `gorm:"column: Introduction; type:VARCHAR(100)  NOT NULL; default:\"\"" json:"introduction"`
	Grade         int    `gorm:"column: Grade; type: TINYINT UNSIGNED NOT NULL;" json:"grade"`
	MailVaild     bool   `gorm:"column: MailVaild; type:BOOLEAN NOT NULL; default:false" json:"mail_vaild"`
	LastLoginTime int64  `gorm:"column: LastLoginTime; type:BIGINT UNSIGNED NOT NULL ;" json:"last_login_time"`
	CreatedAt     int64  `gorm:"column: CreatedAt; type:BIGINT UNSIGNED NOT NULL ;" json:"create_at"`
}

// NotActivatedUser 未進行郵件認證的使用者
// 真正紀錄 Email 的欄位是 NotActEmail
type NotActivatedUser struct {
	User
	NotActEmail     string `gorm:"column: NotActEmail; type:VARCHAR(256) NOT NULL; index:idx_email;" json:"not_act_email"`
	ActaivateToken  string `gorm:"column: EmailActaivateToken; type: VARCHAR(256)  NOT NULL ; uniqueIndex:idx_emat;" json:"email_actaivate_token"`
	EmailExpiration int64  `gorm:"column: EmailExpiration; type:BIGINT UNSIGNED NOT NULL " json:"email_expiration"`
}

// Image datebase image struct
type Image struct {
	ID       int64  `gorm:"column: ID; type: BIGINT UNSIGNED NOT NULL auto_increment; primary_key;" json:"id"`
	HashID   string `gorm:"column: HashID; type: VARCHAR(30);" json:"hash_id"`
	OwnerID  int64  `gorm:"column: OwnerID; type: BIGINT UNSIGNED NOT NULL; index:idx_owner_id;" json:"owner_id"`
	FileName string `gorm:"column: FileName; type: VARCHAR(40);" json:"file_name"`

	Type       string `gorm:"column: Type; type: VARCHAR(10) NOT NULL;" json:"type"`
	Width      int    `gorm:"column: Width; type:INTEGER UNSIGNED NOT NULL " json:"width"`
	Height     int    `gorm:"column: Height; type:INTEGER UNSIGNED NOT NULL " json:"height"`
	Size       int64  `gorm:"column: Size; type:BIGINT UNSIGNED" json:"size"`
	MediumSize int64  `gorm:"column: MediumSize; type:BIGINT UNSIGNED; default:0" json:"medium_size"`

	Title       string `gorm:"column: Title; type:VARCHAR(30) NOT NULL;" json:"title"`
	Description string `gorm:"column: Description; type:VARCHAR(1000) NOT NULL;" json:"description"`
	CreateAt    int64  `gorm:"column: CreatedAt; type:BIGINT UNSIGNED NOT NULL " json:"create_at"`
	UpdateAt    int64  `gorm:"column: UpdateAt; type:BIGINT UNSIGNED NOT NULL " json:"update_at"`
}

// TableName 指定 Setting 表格的名稱
func (Setting) TableName() string {
	return "settings"
}

// TableName 指定 User 表格的名稱
func (User) TableName() string {
	return "users"
}

// TableName 指定 Image 表格的名稱
func (Image) TableName() string {
	return "images"
}

// TableName 指定 NotActivatedUser 表格的名稱
func (NotActivatedUser) TableName() string {
	return "not_activated_user"
}
