package cmd

import (
	"fmt"
	"github.com/douguohai/noticejob/entity"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "push",
	Short: "Ding talk notice",
	Long:  `message push to Ding talk`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("run push...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetMsgFromGit(path string) (entity.DeployMsg, error) {
	// 初始化一个新的仓库对象，指向你的本地Git仓库路径
	r, err := git.PlainOpen(path)
	if err != nil {
		logrus.Infof("读取本地仓库失败 %s", err)
		return entity.DeployMsg{}, err
	}

	// 检查并获取HEAD指向的最新提交
	ref, err := r.Head()
	if err != nil {
		logrus.Infof("检查并获取HEAD指向的最新提交 失败 %s", err)
		return entity.DeployMsg{}, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		logrus.Infof("读取最新提交记录失败 %s", err)
		return entity.DeployMsg{}, err
	}

	// 最后一次提交的提交时间
	commitTime := commit.Author.When

	// 最后一次提交的备注信息（消息）
	message := commit.Message

	//构建消息体
	pushMsg := entity.DeployMsg{ProjectName: "测试项目", PackageVersion: "dev", PackageRunEnv: "打包版本", LastCommitMessage: commitTime.Format(time.DateTime) + " " + message, AccessUrl: "https://www.baidu.com"}

	// 发送钉钉消息
	return pushMsg, nil
}
