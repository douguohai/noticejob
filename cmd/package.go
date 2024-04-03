package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/douguohai/noticejob/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var packageProjectName string    //项目名称
var packagePackageVersion string //版本号
var packagePackageRunEnv string  //打包环境
var packageGitPath string        //git 项目根目录所在
var packageDingToken string      //钉钉 token

var packageMsgContent = "**服务名称**: %s(打包提醒) \n\n **版本类型**: %s\n\n **打包版本**: %s\n\n  **最新提交记录**: \n - %s \n\n  **负载消息**:\n - 流水线打包完成，请在gitlab等待点击部署"

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "package notice",
	Long:  `when ci/cd auto package success ,then push msg to Ding talk`,
	Run: func(cmd *cobra.Command, args []string) {
		if packageDingToken == "" {
			logrus.Error("dingToken is not allow empty")
			return
		}

		pushMsg, err := GetMsgFromGit(packageGitPath)
		if err != nil {
			return
		}
		pushMsg.PackageRunEnv = packagePackageRunEnv
		pushMsg.ProjectName = packageProjectName
		pushMsg.PackageVersion = packagePackageVersion
		SendPackageMessage(pushMsg)
	},
}

func init() {
	packageCmd.Flags().StringVarP(&packageProjectName, "projectName", "n", "", "current project name")
	packageCmd.Flags().StringVarP(&packagePackageVersion, "packageVersion", "v", "", "tag version")
	packageCmd.Flags().StringVarP(&packagePackageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	packageCmd.Flags().StringVarP(&packageGitPath, "gitPath", "g", "./", "project root dir path")
	packageCmd.Flags().StringVarP(&packageDingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(packageCmd)
}

// SendPackageMessage 发送钉钉消息
func SendPackageMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{packageDingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(packageMsgContent, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.LastCommitMessage))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
