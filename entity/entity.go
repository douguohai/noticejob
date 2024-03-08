package entity

// DeployMsg 部署服务消息体
type DeployMsg struct {
	ProjectName       string //项目名称
	PackageVersion    string //版本号
	PackageRunEnv     string //打包环境
	LastCommitMessage string //最后提交消息
	LastCommitTime    string //最后提交时间
	AccessUrl         string //访问地址
	MsgContent        string //消息体
}
