1. 添加系统维护时间段，刷表（修复图像组排序，重置排序）
2. 添加系统维护时间段，刷OSS（没有查询到在 db 里的脏oss数据，删）
3. 清除没有消费的消息队列
4. 查询所有没设置ttl的redis，key，一定有不再需要了能删除的地方
5. fe图像批量接口，分开上传，浏览器压力大
6. ossutil卸载


标准rpc接口描述

+----------------------------------------------------+-------------------------+---------------------------------------------------------+
| service.method                                     | description             | handler                                                 |
+----------------------------------------------------+-------------------------+---------------------------------------------------------+
| SocketService.CreateSocketTicket                   | 创建 WebSocket 连接票据 | socket.(*Service).CreateSocketTicket                    |
| SocketService.ValidateSocketTicket                 | 校验 WebSocket 连接票据 | socket.(*Service).ValidateSocketTicket                  |
| UserService.GetAvatar                              | 获取用户头像            | user.(*Service).GetAvatar                               |
| UserService.UploadAvatar                           | 上传用户自定义头像      | user.(*Service).UploadAvatar                            |
| UserService.DeleteAvatar                           | 删除用户自定义头像      | user.(*Service).DeleteAvatar                            |
| AuthService.Login                                  | 登录                    | auth.(*Service).Login                                   |
| AuthService.Logout                                 | 登出                    | auth.(*Service).Logout                                  |
| AuthService.Me                                     | 获取用户信息            | auth.(*Service).Me                                      |
| AuthService.Refresh                                | 刷新 Token              | auth.(*Service).Refresh                                 |
| AuthService.Register                               | 注册                    | auth.(*Service).Register                                |
| AuthService.ResetPassword                          | 重置密码                | auth.(*Service).ResetPassword                           |
| AuthService.SendEmailCode                          | 发送邮箱验证码          | auth.(*Service).SendEmailCode                           |
| ProjectCoreService.AdvanceProject                  | 项目进度流转            | project/core.(*Service).AdvanceProject                  |
| ProjectCoreService.CreateProject                   | 创建项目                | project/core.(*Service).CreateProject                   |
| ProjectCoreService.DeleteProject                   | 删除项目                | project/core.(*Service).DeleteProject                   |
| ProjectCoreService.ListProjects                    | 获取项目列表            | project/core.(*Service).ListProjects                    |
| ProjectCoreService.CheckProjectAccess              | 校验项目访问权限        | project/core.(*Service).CheckProjectAccess              |
| ProjectProfileService.GetProjectProfile            | 获取项目基础信息        | project/profile.(*Service).GetProjectProfile            |
| ProjectProfileService.GetProjectThumbnail          | 获取项目缩略图          | project/profile.(*Service).GetProjectThumbnail          |
| ProjectProfileService.UploadProjectThumbnail       | 上传项目缩略图          | project/profile.(*Service).UploadProjectThumbnail       |
| ProjectProfileService.DeleteProjectThumbnail       | 删除项目缩略图          | project/profile.(*Service).DeleteProjectThumbnail       |
| ProjectProfileService.UpdateProjectProfile         | 更新项目基础信息        | project/profile.(*Service).UpdateProjectProfile         |
| ProjectAssetsService.GetProjectAssets              | 获取项目图像列表        | project/assets.(*Service).GetProjectAssets              |
| ProjectAssetsService.CreateProjectGroup            | 创建图像组              | project/assets.(*Service).CreateProjectGroup            |
| ProjectAssetsService.DeleteProjectGroup            | 删除图像组              | project/assets.(*Service).DeleteProjectGroup            |
| ProjectAssetsService.MoveProjectGroup              | 移动图像组              | project/assets.(*Service).MoveProjectGroup              |
| ProjectAssetsService.UpdateProjectGroup            | 更新图像组              | project/assets.(*Service).UpdateProjectGroup            |
| ProjectAssetsService.DeleteProjectImage            | 删除图像                | project/assets.(*Service).DeleteProjectImage            |
| ProjectAssetsService.GetProjectImageOriginal       | 获取原图                | project/assets.(*Service).GetProjectImageOriginal       |
| ProjectAssetsService.MoveProjectImage              | 移动图像                | project/assets.(*Service).MoveProjectImage              |
| ProjectAssetsService.ReportProjectImage            | 上报图像                | project/assets.(*Service).ReportProjectImage            |
| ProjectAssetsService.UploadProjectImage            | 上传图像                | project/assets.(*Service).UploadProjectImage            |
| ProjectDetectionService.Ping                       | 智能检测服务探活        | project/detection.(*Service).Ping                       |
| ProjectDetectionService.ReportClassificationResult | 上报图像检测分类结果    | project/detection.(*Service).ReportClassificationResult |
| ProjectDetectionService.ReportReasoningResult      | 上报图像检测推理结果    | project/detection.(*Service).ReportReasoningResult      |
| ProjectDetectionService.ReportSummaryResult        | 上报图像检测总结结果    | project/detection.(*Service).ReportSummaryResult        |
| ProjectReviewService.Ping                          | 人工复核服务探活        | project/review.(*Service).Ping                          |
| ProjectReportService.Ping                          | 评估报告服务探活        | project/report.(*Service).Ping                          |
+----------------------------------------------------+-------------------------+---------------------------------------------------------+

标准http接口描述

+--------+--------------------------------+-----------------------------+---------------------------------------------------+
| method | path                           | description                 | handler                                           |
+--------+--------------------------------+-----------------------------+---------------------------------------------------+
| POST   | /project/assets/group/create   | 创建图像组                  | project/assets.(*Handler).CreateProjectGroup      |
| POST   | /project/assets/group/delete   | 删除图像组                  | project/assets.(*Handler).DeleteProjectGroup      |
| POST   | /project/assets/group/move     | 移动图像组                  | project/assets.(*Handler).MoveProjectGroup        |
| POST   | /project/assets/group/update   | 更新图像组                  | project/assets.(*Handler).UpdateProjectGroup      |
| POST   | /project/assets/image/delete   | 删除图像                    | project/assets.(*Handler).DeleteProjectImage      |
| POST   | /project/assets/image/move     | 移动图像                    | project/assets.(*Handler).MoveProjectImage        |
| POST   | /project/assets/image/report   | 上报图像                    | project/assets.(*Handler).ReportProjectImage      |
| POST   | /project/assets/image/upload   | 上传图像                    | project/assets.(*Handler).UploadProjectImage      |
| POST   | /project/core/advance          | 项目进度流转                | project/core.(*Handler).AdvanceProject            |
| POST   | /project/core/create           | 创建项目                    | project/core.(*Handler).CreateProject             |
| POST   | /project/core/delete           | 删除项目                    | project/core.(*Handler).DeleteProject             |
| POST   | /project/profile/thumbnail     | 上传项目缩略图              | project/profile.(*Handler).UploadProjectThumbnail |
| POST   | /project/profile/update        | 更新项目基础信息            | project/profile.(*Handler).UpdateProjectProfile   |
| POST   | /auth/refresh                  | 刷新 Token                  | auth.(*Handler).Refresh                           |
| POST   | /auth/register                 | 注册                        | auth.(*Handler).Register                          |
| POST   | /auth/reset-password           | 重置密码                    | auth.(*Handler).ResetPassword                     |
| POST   | /auth/login                    | 登录                        | auth.(*Handler).Login                             |
| POST   | /auth/logout                   | 登出                        | auth.(*Handler).Logout                            |
| POST   | /auth/send-email-code          | 发送邮箱验证码              | auth.(*Handler).SendEmailCode                     |
| POST   | /user/avatar                   | 上传用户自定义头像          | user.(*Handler).UploadAvatar                      |
| POST   | /socket/ticket                 | 创建 WebSocket 连接票据     | socket.(*Handler).CreateSocketTicket              |
| GET    | /project/profile/detail        | 获取项目基础信息            | project/profile.(*Handler).GetProjectProfile      |
| GET    | /project/profile/thumbnail     | 获取项目缩略图              | project/profile.(*Handler).GetProjectThumbnail    |
| GET    | /project/assets/list           | 获取项目图像列表            | project/assets.(*Handler).GetProjectAssets        |
| GET    | /project/assets/image/original | 获取原图                    | project/assets.(*Handler).GetProjectImageOriginal |
| GET    | /project/core/list             | 获取项目列表                | project/core.(*Handler).ListProjects              |
| GET    | /auth/me                       | 获取用户信息                | auth.(*Handler).Me                                |
| GET    | /user/avatar                   | 获取用户头像                | user.(*Handler).GetAvatar                         |
| GET    | /socket/setup/assets           | 建立图像资产 WebSocket 连接 | socket.(*Handler).SetupAssetsWebSocket            |
| DELETE | /user/avatar                   | 删除用户自定义头像          | user.(*Handler).DeleteAvatar                      |
| DELETE | /project/profile/thumbnail     | 删除项目缩略图              | project/profile.(*Handler).DeleteProjectThumbnail |
+--------+--------------------------------+-----------------------------+---------------------------------------------------+