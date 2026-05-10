
#### 构建项目

我们推荐使用命令行构建项目，也支持通过 VS Code 的 LaTeX Workshop 插件构建。

##### 通过命令行

###### Makefile (Linux/macOS)

```shell
make all                # compile main.pdf
make ENGINE=$ENGINE all # use $ENGINE (where $ENGINE=-xelatex or -lualatex) to compile main.pdf
make clean              # rm intermediate files
make cleanall           # rm all intermediate files (including .pdf)
make wordcount          # wordcount
```


### 模板配置

#### 文档类选项

在 `main.tex` 的 `\documentclass` 中配置：

```latex
\documentclass[
  oneside,              % 单面打印（默认），使用 twoside 可启用双面打印
  degree=bachelor,      % 学位类型：bachelor（默认），master/doctor 留作扩展
  field=science,        % 专业类别：science 理工科（默认）/ humanities 文科
  fullwidthstop=circle, % 句号样式：circle 保留"。"（默认）/ dot 替换为"．"
  fontset=fandol,       % 字体集，传递给 ctex，默认为 fandol
  times=false,          % true：使用系统 Times New Roman；false：使用 newtx（默认）
  minted=true,          % true：minted 代码高亮（需 Python+Pygments）；false：listings
  biblatex=true,        % true：biblatex+biber（默认）；false：bibtex+gbt7714
]{tongjithesis}

\tjbibresource{bib/note.bib}  % 指定参考文献数据库文件（支持多文件，逗号分隔）
```

### 字体选择

- **Windows 用户**：可直接使用 `fontset=windows`，系统自带 SimSun / SimHei / KaiTi / FangSong 等字体，覆盖更广。
- **跨平台用户**：推荐默认的 `fontset=fandol`（随 TeX Live 安装，零配置）。如需更广字符覆盖，可从 [cjk-fonts-for-ctex](https://github.com/TJ-CSCCG/cjk-fonts-for-ctex) 下载 `adobe` / `founder` / `windows` 等字体并安装到系统，然后切换 `fontset`。

> [!NOTE]
> 安装新字体后请运行 `fc-cache -fv` 刷新字体缓存。

### 代码高亮

1. **`minted`**（默认）：基于 Python Pygments，语法高亮更丰富。需安装 Python 并确保 `pygments` 可用（`pip install pygments`）。
2. **`listings`**：纯 LaTeX 实现，无外部依赖。

在 `main.tex` 中设置 `minted=false` 即可切换。遇到 `minted` 相关错误时，改为 `minted=false` 即可。
