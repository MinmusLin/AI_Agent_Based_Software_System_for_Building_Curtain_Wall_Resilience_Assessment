\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{tongjithesis}[2026/05/18 v1.4.2 Tongji University Thesis Class]

\RequirePackage{kvoptions}
\SetupKeyvalOptions{
  family=tongjithesis,
  prefix=tongjithesis@,
  setkeys=\kvsetkeys
}

\DeclareStringOption[fandol]{fontset}
\DeclareStringOption[circle]{fullwidthstop}
\DeclareBoolOption[false]{times}
\DeclareBoolOption[false]{minted}
\DeclareBoolOption[true]{biblatex}
\DeclareStringOption[bachelor]{degree}
\DeclareStringOption[science]{field}
\DeclareStringOption[algpseudocode]{algo}

\DeclareDefaultOption{\PassOptionsToClass{\CurrentOption}{ctexbook}}
\ProcessKeyvalOptions*

\RequirePackage{etoolbox}
\newif\iftongjithesis@bachelor
\expandafter\ifstrequal\expandafter{\tongjithesis@degree}{bachelor}{%
  \tongjithesis@bachelortrue
}{%
  \ClassWarning{tongjithesis}{type=\tongjithesis@degree\space is reserved for future use. Falling back to bachelor format}%
  \tongjithesis@bachelortrue
}
\newif\iftongjithesis@humanities
\expandafter\ifstrequal\expandafter{\tongjithesis@field}{humanities}{%
  \tongjithesis@humanitiestrue
}{%
  \expandafter\ifstrequal\expandafter{\tongjithesis@field}{science}{}{%
    \ClassWarning{tongjithesis}{field=\tongjithesis@field\space is unknown. Falling back to science format}%
  }%
}
\newif\iftongjithesis@algorithmtwoe
\expandafter\ifstrequal\expandafter{\tongjithesis@algo}{algorithm2e}{%
  \tongjithesis@algorithmtwoetrue
}{}{}
\let\iftjalgorithmtwoe\iftongjithesis@algorithmtwoe

\LoadClass[UTF8,a4paper,zihao=-4,fontset=none]{ctexbook}%

\IfFileExists{tongji-cjk-font-\tongjithesis@fontset.def}
  {\input{tongji-cjk-font-\tongjithesis@fontset.def}}
  {\ClassError{tongjithesis}{Unknown fontset: \tongjithesis@fontset}{}}

\newcommand{\tjfonttitle}{\zihao{-2}}%
\newcommand{\tjfontchapter}{\zihao{-3}}%
\newcommand{\tjfontheading}{\zihao{4}}%
\newcommand{\tjfontbody}{\zihao{-4}}%
\newcommand{\tjfontcaption}{\zihao{5}}%
\newcommand{\tjfontcover}{\zihao{-3}}%
\newcommand{\tjfonttable}{\zihao{-5}}%
\newcommand{\tjfontinfotitle}{\zihao{3}}%
\newcommand{\tjskipbeforesec}{0.6em}%
\newcommand{\tjskipaftersec}{0.6em}%
\newcommand{\tjlinespread}{1.625}%
\newcommand{\tjtoclinespread}{1.625}%
\newcommand{\tjinfofieldspread}{1.6667}%
\newcommand{\tjinfoabstractspread}{1.25}%
\newcommand{\tjinfotitlevspace}{0.8em}%
\newcommand{\tjcoverlinestretch}{0.8}%
\iftongjithesis@humanities
  \newcommand{\tjcovertext}{毕业论文（设计）}
\else
  \newcommand{\tjcovertext}{毕业设计（论文）}
\fi
\iftongjithesis@humanities
  \newcommand{\tjheadertext}{毕业论文（设计）}
\else
  \newcommand{\tjheadertext}{毕业设计（论文）}
\fi

\setlength{\headheight}{43pt}
\RequirePackage[a4paper,top=4.0cm,bottom=2.7cm,left=3.3cm,right=1.8cm]{geometry}
\RequirePackage{fancyhdr}
\RequirePackage{lastpage}
\RequirePackage{eso-pic}
\RequirePackage{emptypage}

\RequirePackage[T1]{fontenc}
\RequirePackage{setspace}
\RequirePackage[inline]{enumitem}
\RequirePackage{alphalph}
\RequirePackage{circledsteps}
\RequirePackage[normalem]{ulem}
\RequirePackage{calc}
\RequirePackage{ifthen}
\RequirePackage{hologo}
\RequirePackage{graphicx}

\RequirePackage[titles]{tocloft}

\RequirePackage{amsmath,amsthm,amsfonts,amssymb,amscd}
\iftongjithesis@times
  \RequirePackage{fontspec}%
  \setmainfont{Times New Roman}%
\else
  \RequirePackage{newtx}%
  \defaultfontfeatures{Ligatures = TeX}%
\fi

\RequirePackage{upgreek}

\RequirePackage{float}
\RequirePackage{multirow}
\RequirePackage{array}
\RequirePackage{booktabs}
\RequirePackage{longtable}
\RequirePackage{threeparttable}
\RequirePackage{threeparttablex}
\RequirePackage{caption}
\DeclareCaptionLabelSeparator{tj@twospace}{\ ~}
\DeclareCaptionFont{tj@songti}{\songti}
\DeclareCaptionFont{tj@caption}{\tjfontcaption}
\captionsetup{labelsep=tj@twospace,font={tj@caption,tj@songti}}
\captionsetup[lstlisting]{skip=12pt}
\captionsetup[listing]{skip=12pt}

\RequirePackage{subcaption}
\captionsetup[sub]{font={tj@caption,tj@songti}}

\iftongjithesis@algorithmtwoe
  \RequirePackage[ruled,linesnumbered,algochapter]{algorithm2e}
  \let\c@algorithm\c@algocf
  \renewcommand{\thealgocf}{\thealgorithm}
\else
  \RequirePackage{algorithm}
  \RequirePackage{algorithmicx}
  \RequirePackage{algpseudocode}
\fi

\RequirePackage{listings}
\definecolor{bg}{rgb}{0.96,0.96,0.96}
\definecolor{commentcolor}{rgb}{0.24,0.48,0.48}
\definecolor{stringcolor}{rgb}{0.73,0.13,0.13}
\definecolor{keywordcolor}{rgb}{0.0,0.502,0.0}
\lstset{
  backgroundcolor=\color{bg},
  basicstyle=\small\ttfamily\linespread{1.0}\selectfont,
  keywordstyle=\color{keywordcolor}\bfseries,
  commentstyle=\color{commentcolor}\itshape,
  stringstyle=\color{stringcolor},
  numbers=left,
  numberstyle=\small\color{gray},
  numbersep=12pt,
  xleftmargin=2.5em,
  xrightmargin=2mm,%
  frame=single,
  rulecolor=\color{gray!40},
  framesep=2mm,
  aboveskip=\baselineskip,
  belowskip=\smallskipamount,
  breaklines=true,
  showstringspaces=false,
  captionpos=b,
  tabsize=4,
  mathescape=true,
  escapeinside=||,
  belowskip=0pt,
}
\renewcommand{\lstlistingname}{代码}
\renewcommand{\lstlistlistingname}{代码列表}

\RequirePackage{fvextra}
\fvset{
  frame=single,
  framesep=2mm,
  rulecolor=\color{gray!40},
  numbers=left,
  numbersep=12pt,
  xleftmargin=2.5em,
  fontsize=\small,
  baselinestretch=1,
  tabsize=4,
  breaklines=true,
  mathescape=true
}
\@ifpackagelater{fvextra}{2024/09/14}{%
  \fvset{bgcolor=bg}%
}{}

\iftongjithesis@minted
  \RequirePackage[chapter]{minted}
  \setminted{linenos,
    frame=single,
    framesep=2mm,
    rulecolor=\color{gray!40},
    fontsize=\small,
    breaklines,
    mathescape,
    escapeinside=||}
  \@ifpackagelater{fvextra}{2024/09/14}{%
    \setminted{bgcolor=bg}%
  }{}
  \renewcommand{\listingscaption}{代码}

\else
  \newfloat{listing}{htbp}{lop}
  \floatname{listing}{代码}

  \newenvironment{minted}[2][]{%
    \begin{lstlisting}[language=#2,#1]%
  }{%
    \end{lstlisting}%
  }

  \newcommand{\inputminted}[3][]{\lstinputlisting[language=#2,#1]{#3}}
\fi

\newenvironment{longlisting}
  {\captionsetup{type=listing,skip=2pt}\lstset{belowskip=10pt}}
  {\medskip}

\RequirePackage{tcolorbox}
\RequirePackage{verbatim}%

\RequirePackage[colorlinks,allcolors=black,unicode=true]{hyperref}
\RequirePackage{nameref}
\RequirePackage[nameinlink,noabbrev]{cleveref}

\RequirePackage{siunitx}[=v2]
\RequirePackage{pgf}
\RequirePackage{tikz}
\usetikzlibrary{shapes.geometric}

\iftongjithesis@biblatex
  \RequirePackage[
    backend=biber,
    style=gb7714-2015,
    natbib=true,
  ]{biblatex}
  \defbibenvironment{bibliography}
    {\list
       {\printtext[labelnumberwidth]{%
          \printfield{labelprefix}%
          \printfield{labelnumber}}}
       {\setlength{\labelwidth}{\labelnumberwidth}%
        \ifdim\labelwidth<0.74cm \setlength{\labelwidth}{0.74cm}\fi
        \setlength{\labelsep}{0pt}%
        \setlength{\leftmargin}{\labelwidth}%
        \setlength{\itemindent}{0pt}%
        \setlength{\itemsep}{\bibitemsep}%
        \setlength{\parsep}{\bibparsep}%
        \renewcommand*{\makelabel}[1]{##1\hfill}}}
    {\endlist}
    {\item}
\else
  \RequirePackage[sort&compress]{gbt7714}
  \bibliographystyle{gbt7714-numerical}
  \setlength{\bibsep}{3pt}
  \newcommand{\tj@bibsetup}[1]{%
    \settowidth\labelwidth{\@biblabel{#1}}%
    \ifdim\labelwidth<0.74cm \setlength{\labelwidth}{0.74cm}\fi
    \setlength{\leftmargin}{\labelwidth}%
    \setlength{\labelsep}{0pt}%
    \setlength{\itemsep}{\bibsep}%
    \setlength{\parsep}{\z@}}
  \AtBeginDocument{\let\@bibsetup\tj@bibsetup}
\fi

\newcommand{\tj@bibfiles}{}
\iftongjithesis@biblatex
  \newcommand{\tjbibresource}[1]{%
    \def\do##1{\addbibresource{##1}}%
    \docsvlist{#1}%
  }
\else
  \newcommand{\tjbibresource}[1]{%
    \def\tj@bibfiles{#1}%
  }
\fi

\iftongjithesis@biblatex
  \newcommand{\makereferences}{%
    \printbibliography[heading=bibintoc,title=参考文献]%
  }
\else
  \newcommand{\makereferences}{%
    \addcontentsline{toc}{chapter}{参考文献}%
    \bibliography{\tj@bibfiles}%
  }
\fi

\setstretch{\tjlinespread}%

\newlength{\bindwidth}
\setlength{\bindwidth}{1em}
\newsavebox{\tjbindpipe}
\newsavebox{\tjbindZ}
\newsavebox{\tjbindD}
\newsavebox{\tjbindX}
\AtBeginDocument{%
  \sbox{\tjbindpipe}{\makebox[\bindwidth][c]{\rule{0.4pt}{3mm}}}%
  \sbox{\tjbindZ}{\makebox[\bindwidth][c]{\zihao{5}装}}%
  \sbox{\tjbindD}{\makebox[\bindwidth][c]{\zihao{5}订}}%
  \sbox{\tjbindX}{\makebox[\bindwidth][c]{\zihao{5}线}}%
}

\newif\iftjbindingline
\tjbindinglinetrue

\newcommand{\tj@bindingcolumn}[1]{%
  \put(\LenToUnit{#1},\LenToUnit{23.60cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{23.12cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{22.65cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{22.16cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{21.68cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{21.20cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{20.72cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{20.24cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{19.76cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{19.29cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{18.80cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{18.32cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{17.84cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{17.37cm}){\usebox{\tjbindZ}}%
  \put(\LenToUnit{#1},\LenToUnit{16.88cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{16.40cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{15.92cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{15.44cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{14.96cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{14.50cm}){\usebox{\tjbindD}}%
  \put(\LenToUnit{#1},\LenToUnit{14.00cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{13.52cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{13.04cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{12.56cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{12.07cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{11.60cm}){\usebox{\tjbindX}}%
  \put(\LenToUnit{#1},\LenToUnit{11.12cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{10.63cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{10.15cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{9.68cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{9.20cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{8.71cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{8.24cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{7.76cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{7.27cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{6.79cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{6.32cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{5.84cm}){\usebox{\tjbindpipe}}%
  \put(\LenToUnit{#1},\LenToUnit{5.35cm}){\usebox{\tjbindpipe}}%
}

\@ifclasswith{tongjithesis}{twoside}{%
  \AddToShipoutPictureBG{%
    \iftjbindingline
      \AtPageLowerLeft{%
        \ifodd\value{page}%
          \tj@bindingcolumn{1.6cm}%
        \else
          \tj@bindingcolumn{19.4cm}%
        \fi
      }%
    \fi
  }%
}{%
  \AddToShipoutPictureBG{%
    \iftjbindingline
      \AtPageLowerLeft{%
        \tj@bindingcolumn{1.6cm}%
      }%
    \fi
  }%
}

\newcommand{\tjdoubleheadrule}{%
  \hrule\@height0.5pt\@width\headwidth
  \vskip1pt
  \hrule\@height0.5pt\@width\headwidth}

\fancypagestyle{firststyle}{
  \fancyhf{}
  \fancyhead[LO]{\hspace*{1cm}\includegraphics[height=1.14cm]{figures/tongji_header.pdf}~}
  \fancyhead[RE]{~\includegraphics[height=1.14cm]{figures/tongji_header.pdf}\hspace*{1cm}}
  \fancyhead[RO,LE]{\tjfontbody ~\tjheadertext~\\}
  \fancyfoot[C]{\tjfontbody \thepage}
  \renewcommand{\headrule}{\tjdoubleheadrule}
  \renewcommand{\footrulewidth}{0pt}
}

\fancypagestyle{mainstyle}{
  \fancyhf{}
  \fancyhead[LO]{\hspace*{1cm}\includegraphics[height=1.14cm]{figures/tongji_header.pdf}~}
  \fancyhead[RE]{~\includegraphics[height=1.14cm]{figures/tongji_header.pdf}\hspace*{1cm}}
  \fancyhead[RO,LE]{\tjfontbody ~\tjheadertext~\\}
  \fancyfoot[RO]{{\tjfontbody 共\quad \pageref{LastPage}\quad 页\quad 第\quad \thepage \quad 页}}
  \fancyfoot[LE]{{\tjfontbody 第\quad \thepage \quad 页\quad 共\quad \pageref{LastPage}\quad 页}}
  \renewcommand{\headrule}{\tjdoubleheadrule}
  \renewcommand{\footrulewidth}{1.8pt}
}

\let\tj@origfrontmatter\frontmatter
\renewcommand{\frontmatter}{%
  \tj@origfrontmatter
  \pagenumbering{Roman}%
  \pagestyle{firststyle}%
  \fancypagestyle{plain}{\pagestyle{firststyle}}%
}
\let\tj@origmainmatter\mainmatter
\renewcommand{\mainmatter}{%
  \tj@origmainmatter
  \pagestyle{mainstyle}%
  \fancypagestyle{plain}{\pagestyle{mainstyle}}%
}
\let\tj@origappendix\appendix
\renewcommand{\appendix}{%
  \tj@origappendix
  \ctexset{chapter/format={\centering\heiti\tjfontheading}}%
  \chapter*{附录}%
  \addcontentsline{toc}{chapter}{附录}%
  \setcounter{section}{0}%
  \renewcommand{\thesection}{\AlphAlph{\value{section}}}%
  \renewcommand{\thesubsection}{\thesection.\arabic{subsection}}%
  \counterwithin{figure}{section}%
  \counterwithin{table}{section}%
  \counterwithin{equation}{section}%
  \counterwithin{algorithm}{section}%
  \counterwithin{listing}{section}%
  \counterwithin{lstlisting}{section}%
  \counterwithin{thm}{section}%
  \counterwithin{cor}{section}%
  \counterwithin{lem}{section}%
  \counterwithin{prop}{section}%
  \counterwithin{conj}{section}%
  \counterwithin{assume}{section}%
  \counterwithin{dfn}{section}%
  \counterwithin{exmp}{section}%
  \counterwithin{rem}{section}%
  \iftongjithesis@humanities
    \ctexset{section/number={\thesection},section/aftername={\hspace{0.5em}}}%
  \else
    \ctexset{section/name={,},section/number={\thesection},section/aftername={.\hspace{0.5em}}}%
  \fi
}
\newcommand{\appendixsection}{\section}
\newcommand{\appendixsubsection}{\subsection}
\let\tj@origbackmatter\backmatter
\renewcommand{\backmatter}{%
  \tj@origbackmatter
  \ctexset{chapter/format={\centering\heiti\tjfontheading}}%
}

\AtBeginEnvironment{table}{\tjfonttable}
\AtBeginEnvironment{tabular}{\tjfonttable}
\AtBeginEnvironment{tabularx}{\tjfonttable}
\AtBeginEnvironment{longtable}{\tjfonttable}
\AtBeginEnvironment{tablenotes}{\tjfonttable}

\newcommand{\tjlongtablecont}[1]{%
  \multicolumn{#1}{r@{\hspace{2em}}}{\zihao{5}\songti 续表~\thetable}\\%
}

\setlistdepth{2}

\setlist[enumerate]{nosep,labelsep=0pt}
\setlist[enumerate,1]{label=（\arabic*）,left=2em,nosep,labelsep=0pt,wide=\parindent}
\setlist[enumerate,2]{label=\Circled{\arabic*},nosep,labelsep=0pt,labelwidth=1em,listparindent=0pt}

\setlist[itemize]{nosep}
\setlist[itemize,1]{wide=\parindent}
\setlist[itemize,2]{wide=2\parindent}

\renewcommand{\cftdotsep}{1}
\setlength{\cftchapindent}{0em}
\setlength{\cftsecindent}{2em}
\setlength{\cftsubsecindent}{4em}
\setlength{\cftsecnumwidth}{1.5em}
\setlength{\cftsubsecnumwidth}{2em}
\setlength{\cftsubsubsecnumwidth}{2.5em}
\setlength{\cftbeforechapskip}{0pt}
\setlength{\cftbeforesecskip}{0pt}
\setlength{\cftbeforesubsecskip}{0pt}
\renewcommand{\contentsname}{目\hspace{0.5em}录}
\renewcommand{\cfttoctitlefont}{\hfill\heiti\tjfontheading}
\renewcommand{\cftaftertoctitle}{\hfill}
\patchcmd{\tableofcontents}{\@starttoc{toc}}{\vspace{\baselineskip}\@starttoc{toc}}%
  {}{\ClassWarning{tongjithesis}{Failed to patch \string\tableofcontents\space for title spacing}}
\renewcommand{\cftchapfont}{\songti\tjfontbody}
\renewcommand{\cftsecfont}{\songti\tjfontbody}
\renewcommand{\cftsubsecfont}{\songti\tjfontbody}
\renewcommand{\cftchapleader}{\cftdotfill{\cftchapdotsep}}
\renewcommand{\cftchapdotsep}{\cftdotsep}
\renewcommand{\cftsecleader}{\cftdotfill{\cftsecdotsep}}
\renewcommand{\cftsecdotsep}{\cftdotsep}
\renewcommand{\cftchappagefont}{\tjfontbody}
\renewcommand{\cftsecpagefont}{\tjfontbody}
\renewcommand{\cftsubsecpagefont}{\tjfontbody}

\let\tj@origtableofcontents\tableofcontents
\renewcommand{\tableofcontents}{%
  \begin{spacing}{\tjtoclinespread}\tj@origtableofcontents\end{spacing}%
}

\RequirePackage{chngcntr}
\counterwithin{figure}{chapter}
\counterwithin{table}{chapter}
\counterwithin{equation}{chapter}
\counterwithin{algorithm}{chapter}
\AtBeginDocument{%
  \counterwithin{listing}{chapter}%
  \counterwithin{lstlisting}{chapter}%
  \let\c@lstlisting\c@listing
}

\ctexset{
  chapter={%
    format={\centering\heiti\tjfontchapter},
    nameformat={},
    name={},
    number={\arabic{chapter}},
    aftername={\quad},
    break={\cleardoublepage},%
    beforeskip={-16pt},
    afterskip={\tjskipaftersec},
  },
  section={%
    format={\heiti\tjfontbody},
    aftername={\hspace{0.5em}},
    beforeskip={\tjskipbeforesec},
    afterskip={\tjskipaftersec},
  },
  subsection={%
    format={\heiti\tjfontbody},
    aftername={\hspace{0.5em}},
    beforeskip={\tjskipbeforesec},
    afterskip={\tjskipaftersec},
    indent={2em},
  },
  subsubsection={%
    format={\heiti\tjfontbody},
    numbering=true,
    number={\AlphAlph{\value{subsubsection}}},
    aftername={.\hspace{0.5em}},
    beforeskip={\tjskipbeforesec},
    afterskip={\tjskipaftersec},
    indent={2em},
    aftertitle={\par},
  },
  paragraph={%
    format={\heiti\tjfontbody},
    numbering=true,
    number={\alphalph{\value{paragraph}}},
    aftername={.\hspace{0.5em}},
    beforeskip={\tjskipbeforesec},
    afterskip={\tjskipaftersec},
    indent={2em},
    aftertitle={\par},
    runin=false,
  },
  appendix/name={},
}

\setcounter{secnumdepth}{5}

\iftongjithesis@humanities
  \renewcommand{\thechapter}{\chinese{chapter}}
  \renewcommand{\thesubsection}{\arabic{subsection}}
  \renewcommand{\thesubsubsection}{\arabic{subsubsection}}
  \ctexset{
    chapter/number    = {\chinese{chapter}},
    chapter/aftername = {、},
    section/name      = {（,）},
    section/number    = {\chinese{section}},
    subsection/number    = {\arabic{subsection}},
    subsection/aftername = {.\hspace{0.5em}},
    subsubsection/name     = {（,）},
    subsubsection/number   = {\arabic{subsubsection}},
    subsubsection/aftername = {\hspace{0.5em}},
    paragraph/number    = {\Circled{\arabic{paragraph}}},
    paragraph/aftername = {\hspace{0.5em}},
  }
  \renewcommand{\cftchapaftersnum}{\!、}
  \setlength{\cftchapnumwidth}{2.5em}
  \setlength{\cftsecindent}{2.5em}
  \setlength{\cftsecnumwidth}{2.5em}
  \renewcommand{\cftsubsecaftersnum}{.\hspace*{0.5em}}
  \setlength{\cftsubsecnumwidth}{1.5em}
  \counterwithout{figure}{chapter}
  \counterwithout{table}{chapter}
  \counterwithout{equation}{chapter}
  \counterwithout{algorithm}{chapter}
  \AtBeginDocument{%
    \counterwithout{listing}{chapter}%
    \counterwithout{lstlisting}{chapter}%
    \counterwithout{thm}{chapter}%
    \counterwithout{cor}{chapter}%
    \counterwithout{lem}{chapter}%
    \counterwithout{prop}{chapter}%
    \counterwithout{conj}{chapter}%
    \counterwithout{assume}{chapter}%
    \counterwithout{dfn}{chapter}%
    \counterwithout{exmp}{chapter}%
    \counterwithout{rem}{chapter}%
  }%
\fi

\renewcommand{\thefootnote}{\Circled{\arabic{footnote}}}

\floatstyle{plaintop}%
\restylefloat{table}
\floatstyle{plain}
\restylefloat{listing}%
\sloppy

\theoremstyle{definition}
\newtheorem{thm}{定理\hspace{0.05pt}}[chapter]
\newtheorem{cor}{推论\hspace{0.05pt}}[chapter]
\newtheorem{lem}{引理\hspace{0.05pt}}[chapter]
\newtheorem{prop}{命题\hspace{0.05pt}}[chapter]
\newtheorem{conj}{猜想\hspace{0.05pt}}[chapter]
\newtheorem{assume}{假设\hspace{0.05pt}}[chapter]
\newtheorem{dfn}{定义\hspace{0.05pt}}[chapter]
\newtheorem{exmp}{例\hspace{0.05pt}}[chapter]
\newtheorem{rem}{注\hspace{0.05pt}}[chapter]
\newtheorem*{pf}{证明}

\providecommand{\dd}{\mathop{}\!\mathrm{d}}
\providecommand{\ee}{\mathrm{e}}
\providecommand{\ii}{\mathrm{i}}
\providecommand{\jj}{\mathrm{j}}
\DeclareRobustCommand\cs[1]{\texttt{\char`\\#1}}
\providecommand\pkg[1]{\href{https://ctan.org/pkg/#1}{\sffamily#1}}

\iftongjithesis@algorithmtwoe
  \SetAlgorithmName{算法}{算法}{算法列表}
  \SetKwInput{KwIn}{输入}
  \SetKwInput{KwOut}{输出}
  \gdef\algocf@Vspace{\vskip\baselineskip\relax}
  \AtBeginEnvironment{algorithm}{\linespread{1.0}\selectfont}
  \AfterEndEnvironment{algorithm}{\vspace{\baselineskip}}
\else
  \renewcommand{\algorithmicrequire}{\textbf{输入: }}
  \renewcommand{\algorithmicensure}{\textbf{输出: }}
  \AtBeginEnvironment{algorithmic}{\linespread{1.0}\selectfont}
  \renewenvironment{algorithm}
  {
    \vskip\baselineskip
    \refstepcounter{algorithm}%
    \hrule height.8pt depth0pt \kern2pt%
    \renewcommand{\caption}[2][\relax]{%
      {\raggedright\textbf{算法~\thealgorithm} ##2\par}%
      \ifx\relax##1\relax%
        \addcontentsline{loa}{algorithm}{\protect\numberline{\thealgorithm}##2}%
      \else%
        \addcontentsline{loa}{algorithm}{\protect\numberline{\thealgorithm}##1}%
      \fi
      \kern2pt\hrule\kern2pt
    }
  }{
    \kern2pt\hrule\relax
    \vskip\baselineskip
  }%
\fi

\iftongjithesis@humanities
  \newcommand*{\chapref}[1]{第\ref{#1}章}
  \newcommand*{\secref}[1]{第\ref{#1}节}
\else
  \newcommand*{\chapref}[1]{第~\ref{#1}~章~}
  \newcommand*{\secref}[1]{第~\ref{#1}~节~}
\fi
\newcommand*{\algoref}[1]{算法~\ref{#1}~}
\renewcommand*{\eqref}[1]{式~(\ref{#1})~}
\newcommand*{\tabref}[1]{\tablename~\ref{#1}~}
\newcommand*{\figref}[1]{\figurename~\ref{#1}~}

\crefformat{listing}{#2代码~#1#3}
\crefmultiformat{listing}{#2代码~#1#3}{和#2代码~#1#3}{、#2代码~#1#3}{和#2代码~#1#3}
\crefrangeformat{listing}{#3代码~#1#4~至#5代码~#2#6}
\crefformat{lstlisting}{#2代码~#1#3}
\crefmultiformat{lstlisting}{#2代码~#1#3}{和#2代码~#1#3}{、#2代码~#1#3}{和#2代码~#1#3}
\crefrangeformat{lstlisting}{#3代码~#1#4~至#5代码~#2#6}
\newcommand*{\listingref}[1]{\lstlistingname~\ref{#1}~}

\iftongjithesis@humanities
  \crefformat{chapter}{#2第#1章#3}
  \crefmultiformat{chapter}{#2第#1章#3}{和#2第#1章#3}{、#2第#1章#3}{和#2第#1章#3}
  \crefrangeformat{chapter}{#3第#1章#4至#5第#2章#6}
  \crefformat{section}{#2第#1节#3}
  \crefmultiformat{section}{#2第#1节#3}{和#2第#1节#3}{、#2第#1节#3}{和#2第#1节#3}
  \crefrangeformat{section}{#3第#1节#4至#5第#2节#6}
\else
  \crefformat{chapter}{#2第~#1~章#3}
  \crefmultiformat{chapter}{#2第~#1~章#3}{和#2第~#1~章#3}{、#2第~#1~章#3}{和#2第~#1~章#3}
  \crefrangeformat{chapter}{#3第~#1~章#4至#5第~#2~章#6}
  \crefformat{section}{#2第~#1~节#3}
  \crefmultiformat{section}{#2第~#1~节#3}{和#2第~#1~节#3}{、#2第~#1~节#3}{和#2第~#1~节#3}
  \crefrangeformat{section}{#3第~#1~节#4至#5第~#2~节#6}
\fi
\crefformat{subsection}{#2第~#1~小节#3}
\crefmultiformat{subsection}{#2第~#1~小节#3}{和#2第~#1~小节#3}{、#2第~#1~小节#3}{和#2第~#1~小节#3}
\crefrangeformat{subsection}{#3第~#1~小节#4至#5第~#2~小节#6}
\crefformat{equation}{#2式~(#1)#3}
\crefmultiformat{equation}{#2式~(#1)#3}{和#2式~(#1)#3}{、#2式~(#1)#3}{和#2式~(#1)#3}
\crefrangeformat{equation}{#3式~(#1)#4~至#5式~(#2)#6}
\crefformat{figure}{#2图~#1#3}
\crefmultiformat{figure}{#2图~#1#3}{和#2图~#1#3}{、#2图~#1#3}{和#2图~#1#3}
\crefrangeformat{figure}{#3图~#1#4~至#5图~#2#6}
\crefformat{table}{#2表~#1#3}
\crefmultiformat{table}{#2表~#1#3}{和#2表~#1#3}{、#2表~#1#3}{和#2表~#1#3}
\crefrangeformat{table}{#3表~#1#4~至#5表~#2#6}
\crefformat{algorithm}{#2算法~#1#3}
\crefmultiformat{algorithm}{#2算法~#1#3}{和#2算法~#1#3}{、#2算法~#1#3}{和#2算法~#1#3}
\crefrangeformat{algorithm}{#3算法~#1#4~至#5算法~#2#6}
\crefformat{algocf}{#2算法~#1#3}
\crefmultiformat{algocf}{#2算法~#1#3}{和#2算法~#1#3}{、#2算法~#1#3}{和#2算法~#1#3}
\crefrangeformat{algocf}{#3算法~#1#4~至#5算法~#2#6}
\crefformat{FancyVerbLine}{#2第~#1~行#3}
\crefmultiformat{FancyVerbLine}{#2第~#1~行#3}{和#2第~#1~行#3}{、#2第~#1~行#3}{和#2第~#1~行#3}
\crefrangeformat{FancyVerbLine}{#3第~#1~行#4至#5第~#2~行#6}
\crefformat{thm}{#2定理~#1#3}
\crefmultiformat{thm}{#2定理~#1#3}{和#2定理~#1#3}{、#2定理~#1#3}{和#2定理~#1#3}
\crefformat{cor}{#2推论~#1#3}
\crefmultiformat{cor}{#2推论~#1#3}{和#2推论~#1#3}{、#2推论~#1#3}{和#2推论~#1#3}
\crefformat{lem}{#2引理~#1#3}
\crefmultiformat{lem}{#2引理~#1#3}{和#2引理~#1#3}{、#2引理~#1#3}{和#2引理~#1#3}
\crefformat{prop}{#2命题~#1#3}
\crefmultiformat{prop}{#2命题~#1#3}{和#2命题~#1#3}{、#2命题~#1#3}{和#2命题~#1#3}
\crefformat{conj}{#2猜想~#1#3}
\crefformat{assume}{#2假设~#1#3}
\crefformat{dfn}{#2定义~#1#3}
\crefformat{exmp}{#2例~#1#3}
\crefformat{rem}{#2注~#1#3}
\crefformat{pf}{#2证明~#1#3}
\crefformat{proof}{#2证明~#1#3}

\newcommand{\crefpairconjunction}{ 和~}
\newcommand{\crefmiddleconjunction}{、}
\newcommand{\creflastconjunction}{ 和~}
\newcommand{\crefpairgroupconjunction}{ 和}
\newcommand{\crefmiddlegroupconjunction}{、}
\newcommand{\creflastgroupconjunction}{ 和}
\newcommand{\crefrangeconjunction}{ 至~}

\providecommand{\XeTeX}{\hologo{XeTeX}}
\providecommand{\BibTeX}{\hologo{BibTeX}}
\providecommand{\BibLaTeX}{\textsc{Bib}\LaTeX}
\providecommand{\LuaTeX}{\hologo{LuaTeX}}
\DeclareRobustCommand{\TongjiThesis}{%
  \texorpdfstring{%
    T{\scshape ongji}\kern.08em T{\scshape hesis}%
  }{TongjiThesis}%
}

\newcolumntype{M}[1]{>{\centering\arraybackslash}m{#1}}
\let\oldtabular\tabular
\let\endoldtabular\endtabular
\newenvironment{information}[1][h!]{
  \renewcommand{\arraystretch}{1.5}
  \setlength{\tabcolsep}{8pt}
  \tjfontcover\oldtabular[#1]
}{\renewcommand{\arraystretch}{1.2}\endoldtabular
}

\newlength{\covertextwidth}
\newlength{\covermaxwidth}
\newlength{\coverparwidth}

\newcommand{\ulinecentered}[1]{
  \bgroup
  \let\\\relax
  \setlength{\covertextwidth}{\widthof{#1}}
  \egroup
  \setlength{\covermaxwidth}{20em}
  \ifdim\covertextwidth > \covermaxwidth
    \ifdim\covertextwidth > 1.95\covermaxwidth
      \setlength{\coverparwidth}{0.53\covertextwidth}
      \pgfmathsetmacro{\coverscale}{\covermaxwidth / \coverparwidth}
      \uline{%
        \scalebox{\coverscale}[1]{%
          \parbox[b]{\coverparwidth}{%
            \centering
            \renewcommand{\baselinestretch}{\tjcoverlinestretch}\selectfont
            #1%
          }%
        }%
      }
    \else
      \uline{%
        \parbox[b]{\covermaxwidth}{%
          \centering
          \renewcommand{\baselinestretch}{\tjcoverlinestretch}\selectfont
          #1%
        }%
      }
    \fi
  \else
    \uline{%
      \parbox[b]{\covermaxwidth}{%
        \centering
        \renewcommand{\baselinestretch}{\tjcoverlinestretch}\selectfont
        #1%
      }%
    }
  \fi
}

\newcommand{\school}[1]{\def\tongjischool{#1}}
\newcommand{\major}[1]{\def\tongjimajor{#1}}
\newcommand{\student}[2]{\def\tongjiauthornumber{#1}\def\tongjiauthor{#2}}
\newcommand{\thesistitle}[2]{\def\tongjithesistitle{#1}\def\tongjithesissubtitle{#2}}
\newcommand{\thesistitleeng}[2]{\def\tongjithesistitleeng{#1}\def\tongjithesissubtitleeng{#2}}
\newcommand{\thesisadvisor}[1]{\def\tongjithesisadvisor{#1}}
\newcommand{\thesisdate}[3]{\def\tongjithesisyear{#1}\def\tongjithesismonth{#2}\def\tongjithesisday{#3}}

\newcommand{\MakeCover}{
  \begin{titlepage}
    \begin{center}
      \global\tjbindinglinefalse
      \pagestyle{empty}%
      \vfill

      \includegraphics[height=2.5cm]{figures/tongji.pdf}\\

      \tjfonttitle\MakeUppercase{\textbf{\tongjiuniversityeng}}\\

      \zihao{-0}\heiti{\tjcovertext}\\
      \vspace{60pt}

      \begin{information}{M{5em}M{22em}}
        \tjfonttitle\heiti{课题名称} & \ulinecentered{\tjfonttitle\heiti{\tongjithesistitle}} \\
        \tjfonttitle\heiti{副\enskip{}标\enskip{}题} & \ulinecentered{\tjfonttitle\heiti{\tongjithesissubtitle}} \\
        \tjfonttitle\heiti{学\quad{}\quad{}院} & \ulinecentered{\tjfonttitle\heiti{\tongjischool}} \\
        \tjfonttitle\heiti{专\quad{}\quad{}业} & \ulinecentered{\tjfonttitle\heiti{\tongjimajor}} \\
        \tjfonttitle\heiti{学生姓名} & \ulinecentered{\tjfonttitle\heiti{\tongjiauthor}} \\
        \tjfonttitle\heiti{学\quad{}\quad{}号} & \ulinecentered{\tjfonttitle\heiti{\tongjiauthornumber}} \\
        \tjfonttitle\heiti{指导教师} & \ulinecentered{\tjfonttitle\heiti{\tongjithesisadvisor}} \\
        \tjfonttitle\heiti{日\quad{}\quad{}期} & \ulinecentered{\tjfonttitle\heiti{\tongjithesisyear{}年\tongjithesismonth{}月\tongjithesisday{}日}} \\
      \end{information}

      \vfill
    \end{center}
  \end{titlepage}
  \global\tjbindinglinetrue
}

\newcommand{\infotype}[1]{\def\tongjiinfotype{#1}}
\newcommand{\infoabstract}[1]{\def\tongjiinfoabstract{#1}}
\newcommand{\infodrawings}[1]{\def\tongjiinfodrawings{#1}}
\newcommand{\infowordcount}[1]{\def\tongjiinfowordcount{#1}}
\newcommand{\infothesiswords}[1]{\def\tongjiinfothesiswords{#1}}
\newcommand{\infomaterials}[1]{\def\tongjiinfomaterials{#1}}
\newsavebox{\tj@titlebox}
\newlength{\tj@titledp}
\newcommand{\infoline}[2]{%
  \noindent #1：#2\par\vspace{0.5em}
}

\newcommand{\tjchecked}{\rlap{$\checkmark$}$\square$}
\newcommand{\tjunchecked}{$\square$}
\newcommand{\tjcheckbox}[1]{%
  \ifthenelse{\equal{\tongjiinfotype}{#1}\OR
              \(\equal{#1}{design}\AND\equal{\tongjiinfotype}{engineering}\)}%
    {\tjchecked}{\tjunchecked}%
}
\newcommand{\ulinebox}[2][c]{%
  \uline{\makebox[#2][#1]{\phantom{X}}}%
}

\newcommand{\MakeInfoPage}{%
  \begin{titlepage}
    \global\tjbindinglinefalse
    \pagestyle{empty}%

    \begin{center}
      \tjfontinfotitle\heiti 同济大学本科\tjheadertext 信息说明页
    \end{center}

    \vspace{1em}

    {\tjfontheading\songti\linespread{\tjinfofieldspread}\selectfont%
    \setlength{\parindent}{0pt}

    \savebox{\tj@titlebox}{\parbox[t]{\dimexpr\linewidth-5em\relax}{%
      \tongjithesistitle
      \ifx\tongjithesissubtitle\empty\else ——\tongjithesissubtitle\fi
    }}%
    \infoline{课题名称}{\usebox{\tj@titlebox}}%
    \settodepth{\tj@titledp}{\usebox{\tj@titlebox}}%
    \addtolength{\tj@titledp}{-\dp\strutbox}%
    \ifdim\tj@titledp>0pt\vspace{\tjinfotitlevspace}\fi

    \infoline{成果类型}{\tjcheckbox{design}\,毕业设计\quad \tjcheckbox{thesis}\,毕业论文}%

    \infoline{学科专业}{\tongjimajor}%

    \infoline{内容简述（请用300字以内简要概述）}{%
      \par\vspace{0.2em}
      {\tjfontbody\linespread{\tjinfoabstractspread}\selectfont%
       \setlength{\parindent}{2em}\tongjiinfoabstract\par}
    }

    \ifthenelse{\equal{\tongjiinfotype}{design}\OR\equal{\tongjiinfotype}{engineering}}{%
      \infoline{毕业设计}{主要图纸~\ifx\tongjiinfodrawings\empty \ulinebox{4em} \else \underline{\makebox[4em][c]{\tongjiinfodrawings}} \fi~张，正文总字数~\ifx\tongjiinfowordcount\empty \ulinebox{6em} \else \underline{\makebox[6em][c]{\tongjiinfowordcount}} \fi~字}%
    }{%
      \infoline{毕业论文}{正文总字数~\ifx\tongjiinfothesiswords\empty \ulinebox{6em} \else \underline{\makebox[6em][c]{\tongjiinfothesiswords}} \fi~字}%
    }%

    \infoline{随附资料}{%
      \ifx\tongjiinfomaterials\empty
        \par\vspace{0.5em}
      \else
        \begin{enumerate}[label=\arabic*.,leftmargin=2em,labelsep=0.5em,itemsep=0.2em,parsep=0pt,topsep=0pt] \tongjiinfomaterials \end{enumerate}
      \fi
    }%
    （请列出所有提交的支撑材料名称，如：毕业设计--全套图纸、毕业作品、计算书、程序代码、附录等）
    }
  \end{titlepage}
  \setstretch{\tjlinespread}
  \global\tjbindinglinetrue
}

\newcommand{\MakeAbstract}[2]{
  \begin{center}\bfseries\heiti\tjfonttitle
    \ \\
    \vspace{-0.5em}
    \tongjithesistitle
    \ifx\tongjithesissubtitle\empty\else
      \\[0.2em]{\zihao{3}\tongjithesissubtitle}%
    \fi
  \end{center}
  \vspace{-0.8em}
  \begin{center}\heiti\tjfontheading
    \ifthenelse{\equal{\tongjiinfotype}{engineering}}{设计总说明}{摘\hspace{0.5em}要}
  \end{center}

  \begin{spacing}{\tjlinespread}
    \setlength{\parindent}{2em}
    \songti\tjfontbody #1
    \par
  \end{spacing}

  \songti\tjfontbody
  \vspace{\baselineskip}
  {\noindent\textbf{关键词：}#2}\par
}

\newcommand{\MakeAbstractEng}[2]{
  \clearpage
  \begin{center}\tjfonttitle
    \ \\
    \vspace{-0.5em}
    \textbf{\tongjithesistitleeng}
    \ifx\tongjithesissubtitleeng\empty\else
      \\[0.2em]{\zihao{3}\textbf{\tongjithesissubtitleeng}}%
    \fi
  \end{center}
  \vspace{-0.8em}
  \begin{center}\tjfontheading
    \textbf{ABSTRACT}
  \end{center}

  \begin{spacing}{\tjlinespread}
    \setlength{\parindent}{2em}
    \tjfontbody #1
    \par
  \end{spacing}

  \tjfontbody
  \vspace{\baselineskip}
  {\noindent\textbf{Key words: }#2}\par
}

\newif\iftongjithesis@fullwidth@active
\tongjithesis@fullwidth@activefalse
\expandafter\ifstrequal\expandafter{\tongjithesis@fullwidthstop}{dot}
  {\tongjithesis@fullwidth@activetrue}{}
\expandafter\ifstrequal\expandafter{\tongjithesis@fullwidthstop}{true}
  {\tongjithesis@fullwidth@activetrue}{}
\iftongjithesis@fullwidth@active
  \catcode`\。=\active
  \newcommand{。}{\ifmmode\text{．}\else ．\fi}%
\fi

\AtEndOfPackage{
  \InputIfFileExists{tongjithesis.cfg}{}{}
}
