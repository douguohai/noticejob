package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"noticejob/entity"
)

var PackageMsgTemp = "**服务名称**: %s(打包提醒) \n\n **版本类型**: %s\n\n **打包版本**: %s\n\n  **最新提交记录**: \n - %s \n\n  **负载消息**:\n - 流水线打包完成，请在gitlab等待点击部署"

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "package notice",
	Long:  `when ci/cd auto package success ,then push msg to Ding talk`,
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
		SendPackageMessage(pushMsg)
	},
}

func init() {
	packageCmd.Flags().StringVarP(&projectName, "projectName", "n", "", "current project name")
	packageCmd.Flags().StringVarP(&packageVersion, "packageVersion", "v", "", "tag version")
	packageCmd.Flags().StringVarP(&packageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	packageCmd.Flags().StringVarP(&gitPath, "gitPath", "g", "./", "project root dir path")
	packageCmd.Flags().StringVarP(&dingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(packageCmd)
}

// SendPackageMessage 发送钉钉消息
func SendPackageMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{dingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(PackageMsgTemp, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.LastCommitMessage))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
