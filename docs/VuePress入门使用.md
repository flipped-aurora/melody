# VuePress 搭建文档手册

## 前提

- node 环境
- npm

## 安装

```
npm install -g vuepress
```

## 开始

1. 新建个文件夹
2. 初始化，生成pakcage.json
```
npm init -y
```
3. vscode打开项目
4. 新建`docs`目录
5. `docs下`新建`.vuepress`目录
6. `.vuepress`下新建`config.json`
```config.json
module.exports = {
    title: "Melody Docs",
    description: "Documents of melody"
}
```
7. `docs`下新建`README.md`
8. 尝试启动
```
sudo vuepress dev docs
```
 