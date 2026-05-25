# AI Agent-Based Software System for Building Curtain Wall Resilience Assessment

> [!CAUTION]
> ### 免责声明 | Disclaimer
>
> The code and materials contained in this repository are intended for personal learning and research purposes only and may not be used for any commercial purposes. Other users who download or refer to the content of this repository must strictly adhere to the **principles of academic integrity** and must not use these materials for any form of homework submission or other actions that may violate academic honesty. I am not responsible for any direct or indirect consequences arising from the improper use of the contents of this repository. Please ensure that your actions comply with the regulations of your school or institution, as well as applicable laws and regulations, before using this content. If you have any questions, please contact me via [email](mailto:minmuslin@outlook.com).
>
> 本仓库包含的代码和资料仅用于个人学习和研究目的，不得用于任何商业用途。请其他用户在下载或参考本仓库内容时，严格遵守**学术诚信原则**，不得将这些资料用于任何形式的作业提交或其他可能违反学术诚信的行为。本人对因不恰当使用仓库内容导致的任何直接或间接后果不承担责任。请在使用前务必确保您的行为符合所在学校或机构的规定，以及适用的法律法规。如有任何问题，请通过[电子邮件](mailto:minmuslin@outlook.com)与我联系。

Design and Implementation of an AI Agent-Based Software System for Building Curtain Wall Resilience Assessment.

基于 AI Agent 的建筑幕墙韧性评估软件系统设计与实现。

本毕业设计面向建筑幕墙人工巡检效率低、主观性强、难以处理海量图像等问题，设计并实现了一套基于 AI Agent 的建筑幕墙韧性评估软件系统。系统亮点包括面向建筑幕墙多材质图像的 Agent 自主任务调度方法和 LLM 分层信息整合策略：前者由分类 Agent 根据图像内容自主判断需要执行的原子检测能力，并将开放语义判断约束为工程系统可消费的任务代码；后者将大语言模型从简单文本生成工具扩展为“路由、压缩、汇总”的智能节点，有效解决了海量数据直接输入 LLM 模型带来的上下文过载和重点淹没问题。作者现已完成原型系统的部署和实际场景试用，获得同济大学土木工程学院科研团队认可，为建筑幕墙数字运维提供了一种兼具工程可实施性和智能决策能力的解决方案。

## 仓库组成

* `/.github/workflows`
GitHub CI/CD 工作流配置

* `/assets`
存放 `README.md` 文件所需的相关图片资源

* `/Dataset/Earthquake_Engineering_Hall`
同济大学地震工程馆建筑幕墙数据集

* `/Documentations`
毕业设计（论文）归档文件

* `/Source_Code`
建筑幕墙韧性评估软件系统源代码

* `/Thesis_LaTeX`
毕业设计（论文） $\LaTeX$ 源代码

## 项目概览

<p align="center">
  <img src="assets/Final_Defense_01.png" width="48%">
  <img src="assets/Final_Defense_02.png" width="48%">
  <br>
  <img src="assets/Final_Defense_03.png" width="48%">
  <img src="assets/Final_Defense_04.png" width="48%">
  <br>
  <img src="assets/Final_Defense_05.png" width="48%">
  <img src="assets/Final_Defense_06.png" width="48%">
  <br>
  <img src="assets/Final_Defense_07.png" width="48%">
  <img src="assets/Final_Defense_08.png" width="48%">
  <br>
  <img src="assets/Final_Defense_09.png" width="48%">
  <img src="assets/Final_Defense_10.png" width="48%">
  <br>
  <img src="assets/Final_Defense_11.png" width="48%">
  <img src="assets/Final_Defense_12.png" width="48%">
  <br>
  <img src="assets/Final_Defense_13.png" width="48%">
  <img src="assets/Final_Defense_14.png" width="48%">
  <br>
  <img src="assets/Final_Defense_15.png" width="48%">
  <img src="assets/Final_Defense_16.png" width="48%">
  <br>
  <img src="assets/Final_Defense_17.png" width="48%">
  <img src="assets/Final_Defense_18.png" width="48%">
  <br>
  <img src="assets/Final_Defense_19.png" width="48%">
  <img src="assets/Final_Defense_20.png" width="48%">
  <br>
</p>

## 摘要

建筑幕墙是现代建筑的重要外围护结构，其服役状态直接影响建筑安全和运维质量。现实工程中，幕墙巡检长期依赖人工经验，存在效率低、主观性强、难以处理海量图像等问题。已有计算机视觉方法能够完成裂缝、锈蚀、破损等局部缺陷检测，但多数方案仍停留在单模型、单图像、单缺陷识别层面，只能检测某张图像是否存在某类缺陷，缺少多组巡检图像、多种幕墙材质、多类检测结果之间的汇总融合。对于土木工程领域用户而言，实际需求并不是孤立的单图检测结论，而是面向整栋建筑形成的多维度韧性评估结论。因此，为打通从项目级图像资产管理、单图检测结果到建筑级韧性评估报告的完整链路，本文设计并实现了一套基于 AI Agent 的建筑幕墙韧性评估软件系统。

本文的研究重点是 AI Agent 参与工程任务规划、能力调度与信息整合的系统性方法。本文工作的主要亮点包括两个方面：其一，提出面向建筑幕墙多材质图像的 Agent 自主任务调度方法。该方法由分类 Agent 根据图像内容自主判断需要执行的原子检测能力，并将开放语义判断约束为工程系统可消费的任务代码；其二，提出 LLM 分层信息整合策略。该策略将所有多源检测结果、项目背景、图像分组和人工复核意见等信息进行分层整合，将大语言模型从简单文本生成工具扩展为“路由、压缩、汇总”的智能节点，有效解决了海量数据直接输入 LLM 模型带来的上下文过载和重点淹没问题。

在工程实现上，原型系统主要使用 Golang、Python 和 TypeScript 编程语言、React、Gin 和 PyTorch 框架、gRPC 服务通信协议、RocketMQ 消息队列、WebSocket 实时消息推送、MySQL 关系型数据库、Redis 非关系型数据库、MinIO 对象存储和 Docker 容器化等技术。工程链路为 Agent 提供高可用、高性能、高并发的运行环境，原型系统提供从建筑图像资产构建、Agent 智能检测、人工复核到韧性评估报告的完整闭环工作链。作者现已完成原型系统的部署和实际场景试用，获得同济大学土木工程学院科研团队认可，为建筑幕墙数字运维提供了一种兼具工程可实施性和智能决策能力的解决方案。

**关键词**：AI Agent，建筑幕墙，任务编排，韧性评估

## 谢辞

十八岁刚进入同济大学时，我以为二十二岁会是某种写好的答案，但是真到了这一天，才发现那是一道更复杂、更迷茫、也更没有标准答案的题目。

首先，衷心感谢我的指导教师黄杰老师。感谢黄老师在选题确定、系统设计、论文撰写和整体方向把握上给予的指导与帮助。老师的指导帮助我不断收束主线，使我最终明确了核心技术路线与工作创新点，也让我能够把一个庞杂的软件系统整理成清晰的论文结构。毕业设计能够从最初的想法推进到最终成稿，离不开老师在关键阶段的指导、支持与包容。

感谢沈坚老师。高级语言程序设计是我在本科阶段收获最大的课程，也是我真正走进代码世界的起点。沈老师扎实而严谨的教学，为我打下了 C/C++ 编程基础，也帮助我较早培养起了计算机思维。后来无论是写课程项目、完成毕业设计，还是进入工作后编写需要长期稳定运行、服务真实用户的生产代码，我都能感受到那段程序入门训练留下的影响。那些当时看似只是语法、调试和作业的训练，后来都慢慢变成了工程实践中的底气。

感谢我的家人。感谢你们在我本科四年求学过程中始终给予理解、支持和信任。无论是课程压力、竞赛准备，还是实习求职、本科就业，你们都没有用某一种固定答案要求我，而是尊重我的判断，支持我去寻找更适合自己的道路。

感谢 iGEM 2024 Tongji-Software 的所有队友。那是一段非常特别的经历：我们从最初的想法出发，一起讨论、开发、部署，最后在巴黎斩获国际金奖。竞赛本身带给我的不只是奖项，更是一次跨学科协作、产品构思和团队推进的完整训练。它让我意识到，编程能力只是软件工程的一部分，更重要的是如何把想法组织起来、把问题讲清楚，并最终把系统交付出去。

感谢同济大学。感谢西南十二楼党建室、通达馆和济事楼 319 学生技术俱乐部，为我提供了学习、熬夜和赶项目的场所。许多代码、报告、答辩材料和技术方案，都是在这些地方一点点完成的。

也感谢一直陪伴我的她。大学四年里，有许多焦虑、低落和自我怀疑的时刻，也有许多因为项目跑通、面试通过、论文完成而放松下来的瞬间，谢谢你愿意听我讲那些并不总是有趣的技术问题，见证这些具体而琐碎的时刻。

最后，感谢一直努力的自己。回顾本科四年，我庆幸没有只把大学当成绩点和排名的竞赛场。工程能力的培养不能只依赖课堂学习，更需要长期主动实践、持续自学和真实项目的反复锤炼。正是在课程之外的探索与投入中，我从 C/C++ 写到 Java、Golang、TypeScript，从单机系统写到分布式系统，从简单接口写到完整部署，逐渐建立起自己的技术体系。

曾经我总觉得大学时光还有很长，很多事情以后再做也来得及，真正写到毕业论文最后一页时，才发现四年其实过得很快。这里有让我疲惫的时刻，也有让我被推动着成长的时刻；有遗憾，也有幸运。它并没有给我一条标准答案，而是让我在不断选择中逐渐看清自己。

谨以此文，献给我的本科四年。愿自己在未来仍能保持清醒的自我认知，保持对技术的热情与好奇心，也保持继续学习和解决真实问题的能力。
