# LIVE LIFE 文档索引

这个目录现在分成三类文档：

1. 固定需求与架构：多人协作时先看这些，避免每个人按自己的理解改。
2. 开发与部署流程：说明本地、GitHub、阿里云之间怎么协作。
3. 设计方案记录：V2、V3 以及未来 V4/V5 的视觉方向可以变化，但不应该反向影响后端业务契约。

## 固定需求与整体架构

- [requirements-analysis.md](requirements-analysis.md)
- [architecture-overview.md](architecture-overview.md)
- [database-schema-draft.md](database-schema-draft.md)

用途：

- 固定 Shows / CD 严选 / Archive / Connect 的业务边界。
- 固定不设置顶层 Shop 页面，购买从 CD 严选单品跳外部 Shop。
- 固定三语言、LIVE LIFE 品牌展示、Connect 联系入口。
- 固定当前数据库选型和未来升级方向。

## 后端与 Clean Architecture

- [backend-detailed-design.md](backend-detailed-design.md)
- [backend-clean-architecture-review.md](backend-clean-architecture-review.md)

用途：

- 说明后端分层后的目录、职责和依赖方向。
- 固定后端业务逻辑和 API 契约。
- 说明未来即使前端设计从 V3 改到 V4，后端也不跟着视觉版本变化。

## GitHub 多人协作与部署

- [github-collaboration-workflow.md](github-collaboration-workflow.md)
- [cicd-git-push-deploy.md](cicd-git-push-deploy.md)
- [local-development.md](local-development.md)
- [alicloud-tokyo-p0-deployment.md](alicloud-tokyo-p0-deployment.md)
- [alicloud-tokyo-p0-status.md](alicloud-tokyo-p0-status.md)

用途：

- 说明当前本地直推阿里云工作流和未来 GitHub PR 工作流的区别。
- 说明 main / develop 分支规则。
- 说明 1GB 阿里云服务器不在服务器上构建，而是接收预构建 release 产物。
- 说明 GitHub 建好后，其他人如何 fork/branch/PR 给负责人 review。

## 前端设计方案

- [frontend-design-variants.md](frontend-design-variants.md)
- [product-architecture-and-ui-approval.md](product-architecture-and-ui-approval.md)
- [v3-design-approval.md](v3-design-approval.md)

用途：

- 记录 V2、V3 的设计语言。
- 说明右上角 review 下拉框如何选择不同视觉方案。
- 前端视觉可以继续改，但不能随意改后端 API 契约。

## 后续建议补充

- `content-data-entry-guide.md`：后端管理系统完成前，先规定演出、CD、黑胶、Archive 内容怎么录入。
- `external-link-policy.md`：规定票站、BASE、Instagram、小红书等外部链接如何展示，避免误导为站内交易。
- `brand-language-guide.md`：规定 LIVE LIFE、三语言、英文全大写、乐队名纹理和授权提示边界。
- `future-auth-checkout-notes.md`：等决定是否做登录注册和站内支付后，再整理账户、订单、售后设计。
