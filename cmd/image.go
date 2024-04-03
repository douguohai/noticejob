package cmd

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/douguohai/noticejob/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var imageProjectName string    //项目名称
var imagePackageVersion string //版本号
var imagePackageRunEnv string  //打包环境
var imageAccessUrl string      //访问地址
var imageGitPath string        //git 项目根目录所在
var imageDingToken string      //钉钉 token

var ImageBuildMsgTemp = "**服务名称**: %s(镜像构建成功)\n\n **版本类型**: %s\n\n **打包版本**: %s\n\n   **容器地址**: %s \n\n **最新提交记录**: \n - %s \n\n"

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "just build image notice",
	Long:  `when success update k8s or docker ,then push msg to Ding talk`,
	Run: func(cmd *cobra.Command, args []string) {
		if imageDingToken == "" {
			logrus.Error("dingToken is not allow empty")
			return
		}
		pushMsg, err := GetMsgFromGit(imageGitPath)
		if err != nil {
			return
		}
		pushMsg.PackageRunEnv = imagePackageRunEnv
		pushMsg.ProjectName = imageProjectName
		pushMsg.PackageVersion = imagePackageVersion
		pushMsg.AccessUrl = imageAccessUrl
		SendImageBuildMessage(pushMsg)
	},
}

func init() {
	imageCmd.Flags().StringVarP(&imageProjectName, "projectName", "n", "", "current project name")
	imageCmd.Flags().StringVarP(&imagePackageVersion, "packageVersion", "v", "", "tag version")
	imageCmd.Flags().StringVarP(&imagePackageRunEnv, "packageRunEnv", "e", "", "run env dev|alpha|beta|prod")
	imageCmd.Flags().StringVarP(&imageGitPath, "gitPath", "g", "./", "project root dir path")
	imageCmd.Flags().StringVarP(&imageAccessUrl, "accessUrl", "u", "暂未提供链接", "dockerhub image url ")
	imageCmd.Flags().StringVarP(&imageDingToken, "dingToken", "t", "", "ding talk reboot token")
	rootCmd.AddCommand(imageCmd)
}

// SendImageBuildMessage 发送钉钉消息
func SendImageBuildMessage(pushMsg entity.DeployMsg) {
	// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{imageDingToken}
	cli := dingtalk.InitDingTalk(dingToken, ".")
	err := cli.SendMarkDownMessage("服务更新-更新日志", fmt.Sprintf(ImageBuildMsgTemp, pushMsg.ProjectName, pushMsg.PackageRunEnv, pushMsg.PackageVersion, pushMsg.AccessUrl, pushMsg.LastCommitMessage))
	if err == nil {
		logrus.Infof("推送消息成功")
	} else {
		logrus.Infof("推送消息失败 %s", err)
	}
}
