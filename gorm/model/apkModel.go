package model

import "time"

// 应用资源信息表
type ApkModel struct {
	ApkId          int       `gorm:"column:apk_id" json:"apk_id"`             //应用ID
	ApkName        string    `gorm:"column:apk_name" json:"apk_name"`         //应用名称
	ApkResType     string    `gorm:"column:apk_res_type" json:"apk_res_type"` //应用资源类型
	Hot            int       `gorm:"column:hot" json:"hot"`                   //是否热门（0不是  1是热门）
	DownloadNum    int       `gorm:"column:download_num" json:"download_num"` //下载次数(W）
	VersionCode    int       `gorm:"column:version_code" json:"version_code"` //版本号
	VersionName    string    `gorm:"column:version_name" json:"version_name"` //版本名称
	FileMd5        string    `gorm:"column:file_md5" json:"file_md5"`
	FileSize       string    `gorm:"column:file_size" json:"file_size"`
	FileName       string    `gorm:"column:file_name" json:"file_name"`         //应用文件名
	PackageName    string    `gorm:"column:package_name" json:"package_name"`   //应用包名
	DownloadUrl    string    `gorm:"column:download_url" json:"download_url"`   //应用下载地址
	LocalPath      string    `gorm:"column:local_path" json:"local_path"`       //应用服务器地址
	Status         int       `gorm:"column:status" json:"status"`               //应用状态:0删除；-1下线；1上线
	Mark           string    `gorm:"column:mark" json:"mark"`                   //应用标记
	Company        string    `gorm:"column:company" json:"company"`             //公司信息
	CompanyType    int       `gorm:"column:company_type" json:"company_type"`   //0:第三方 1:天奕达 2:合作方
	IsForcedUp     int       `gorm:"column:is_forced_up" json:"is_forced_up"`   //(未使用)
	Integral       int       `gorm:"column:integral" json:"integral"`           //167版本是否强制显示：（未使用）
	IsSpread       int       `gorm:"column:is_spread" json:"is_spread"`         //solr查询排位因数
	IsTime         int       `gorm:"column:is_time" json:"is_time"`             //是否分时段（未使用）
	IsBeginTime    string    `gorm:"column:is_begin_time" json:"is_begin_time"` //分时段开始时间（未使用）
	IsEndTime      string    `gorm:"column:is_end_time" json:"is_end_time"`     //分时段结束时间（未使用）
	GradeId        int       `gorm:"column:grade_id" json:"grade_id"`           //等级
	Stars          int       `gorm:"column:stars" json:"stars"`                 //星级(0:一星, 1:二星, 2:三星, 3:四星, 4:五星)
	Security       int       `gorm:"column:security" json:"security"`           //安全(0:不安全, 1:已经通过安全检测)（未使用）
	Charge         int       `gorm:"column:charge" json:"charge"`               //资费(0:免费, 1:部分功能付费)
	IsAd           int       `gorm:"column:is_ad" json:"is_ad"`                 //广告(0:无广告, 1:有广告)
	ActivityUrl    string    `gorm:"column:activity_url" json:"activity_url"`   //活动地址
	SearchDownload int       `gorm:"column:search_download" json:"search_download"`
	DownNet        int       `gorm:"column:down_net" json:"down_net"`           //静默下载网络配置(1:2G 2:3G 3:WIFI 4:4G 5:5G),默认为wifi
	SilentUpEqu    int       `gorm:"column:silent_up_equ" json:"silent_up_equ"` //705静默更新并安装
	SilentUpGre    int       `gorm:"column:silent_up_gre" json:"silent_up_gre"` //706静默更新并安装
	DroiTest       int       `gorm:"column:droi_test" json:"droi_test"`
	ProductType    int8      `gorm:"column:product_type" json:"product_type"`   //产品分类(0单机、1网游、2应用、3联运单机、4联运网游、5其他、6内测、7预约网游)
	IpType         int8      `gorm:"column:ip_type" json:"ip_type"`             //IP产品(0:不是IP产品 1:是IP产品)
	SignatureMd5   string    `gorm:"column:signature_md5" json:"signature_md5"` //应用签名MD5,177版本前夕添加,仅当前版本之后的运营干预应用有此数据（新增修改联运）（未使用）
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`     //创建时间
	ModifyTime     time.Time `gorm:"column:modify_time" json:"modify_time"`     //修改时间
	Creator        string    `gorm:"column:creator" json:"creator"`             //创建人
	From           string    `gorm:"column:from" json:"from"`                   //来自哪里
	IsSearch       int8      `gorm:"column:is_search" json:"is_search"`         //是否支持搜索
	Video          string    `gorm:"column:video" json:"video"`                 //视频
}

func (ApkModel) TableName() string {
	return "apk"
}
