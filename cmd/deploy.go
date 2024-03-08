package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/douguohai/noticejob/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DeployMsgTemp = "**服务名称**: %s(部署成功)\n\n **版本类型**: %s\n\n **打包版本**: %s\n\n   **访问地址**: %s \n\n **最新提交记录**: \n - %s \n\n"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy notice",
	Long:  `when success update k8s or docker ,then push msg to Ding talk`,
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
		pushMsg.PackageVersion = packageVersion
		SendDeployMessage(pushMsg)
	},
}

func init() {
	deployCmd.Flags().StringVarP(&projectName, "projectName", "n", "", "current project name")
	deployCmd.Flags().StringVarP(&packageVersion, "packageVersion", "v", "", "tag version")
	deployCmd.Flags().StringVarP(&packageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	deployCmd.Flags().StringVarP(&gitPath, "gitPath", "g", "./", "project root dir path")
	deployCmd.Flags().StringVarP(&accessUrl, "accessUrl", "u", "暂未提供链接", "web access url")
	deployCmd.Flags().StringVarP(&dingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(deployCmd)
}

// SendDeployMessage 发送钉钉消息
func SendDeployMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{dingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(DeployMsgTemp, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.AccessUrl, pushMsg.LastCommitMessage))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
