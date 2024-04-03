package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/douguohai/noticejob/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deployProjectName string    //项目名称
var deployPackageVersion string //版本号
var deployPackageRunEnv string  //打包环境
var deployAccessUrl string      //访问地址
var deployGitPath string        //git 项目根目录所在
var deployDingToken string      //钉钉 token

var DeployMsgTemp = "**服务名称**: %s(部署成功)\n\n **版本类型**: %s\n\n **打包版本**: %s\n\n   **访问地址**: %s \n\n **最新提交记录**: \n - %s \n\n"

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy notice",
	Long:  `when success update k8s or docker ,then push msg to Ding talk`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("deployDingToken", deployDingToken)
		if deployDingToken == "" {
			logrus.Error("dingToken is not allow empty")
			return
		}
		pushMsg, err := GetMsgFromGit(deployGitPath)
		if err != nil {
			return
		}
		pushMsg.PackageRunEnv = deployPackageRunEnv
		pushMsg.ProjectName = deployProjectName
		pushMsg.PackageVersion = deployPackageVersion
		pushMsg.AccessUrl = deployAccessUrl
		SendDeployMessage(pushMsg)
	},
}

func init() {
	deployCmd.Flags().StringVarP(&deployProjectName, "projectName", "n", "", "current project name")
	deployCmd.Flags().StringVarP(&deployPackageVersion, "packageVersion", "v", "", "tag version")
	deployCmd.Flags().StringVarP(&deployPackageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	deployCmd.Flags().StringVarP(&deployGitPath, "gitPath", "g", "./", "project root dir path")
	deployCmd.Flags().StringVarP(&deployAccessUrl, "accessUrl", "u", "暂未提供链接", "web access url")
	deployCmd.Flags().StringVarP(&deployDingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(deployCmd)
}

// SendDeployMessage 发送钉钉消息
func SendDeployMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{deployDingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(DeployMsgTemp, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.AccessUrl, pushMsg.LastCommitMessage))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
