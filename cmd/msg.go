package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/douguohai/noticejob/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var projectName string    //项目名称
var packageVersion string //版本号
var packageRunEnv string  //打包环境
var gitPath string        //git 项目根目录所在
var dingToken string      //钉钉 token
var msgContent string     //普通消息体

var MsgTemp = "**服务名称**: %s(消息提醒)\n\n **版本类型**: %s\n\n **打包版本**: %s\n\n  **消息内容**: \n  %s \n\n"

var msgCmd = &cobra.Command{
	Use:   "public",
	Short: "public message notice",
	Long:  `push public message  to Ding talk`,
	Run: func(cmd *cobra.Command, args []string) {
		if dingToken == "" {
			logrus.Error("dingToken is not allow empty")
			return
		}
		pushMsg, err := GetMsgFromGit(gitPath)
		if err != nil {
			return
		}
		pushMsg.PackageRunEnv = packageRunEnv
		pushMsg.ProjectName = projectName
		pushMsg.PackageVersion = packageVersion
		pushMsg.MsgContent = msgContent
		SendMessage(pushMsg)
	},
}

func init() {
	msgCmd.Flags().StringVarP(&projectName, "projectName", "n", "", "current project name")
	msgCmd.Flags().StringVarP(&packageVersion, "packageVersion", "v", "", "tag version")
	msgCmd.Flags().StringVarP(&packageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	msgCmd.Flags().StringVarP(&msgContent, "msgContent", "m", "", "public text")
	msgCmd.Flags().StringVarP(&dingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(msgCmd)
}

// SendMessage 发送钉钉消息
func SendMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{dingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(MsgTemp, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.MsgContent))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
