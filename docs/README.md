# LIVE LIFE 文档索引

这个目录现在分成两类文档：

1. 长期开发文档：需求、后端、前端设计版本、本地启动、部署。
2. 阶段审批文档：V2、V3 的设计审批稿，用来记录当时为什么这么定方向。

设计视觉以后会继续变化，但核心需求、后端 API、购买路径和联系路径先保持稳定。

## 1. 需求与范围

- [architecture-overview.md](architecture-overview.md)
- [requirements-analysis.md](requirements-analysis.md)

用途：

- 先看整体架构图，再看需求细节。
- 记录 LIVE LIFE 当前产品目标。
- 固定不会随设计版本频繁变化的需求。
- 说明 Shows、CD 严选、Archive、Connect 的边界。
- 说明暂缓登录注册和站内支付的原因。

## 2. 后端详细设计

- [backend-detailed-design.md](backend-detailed-design.md)
- [database-schema-draft.md](database-schema-draft.md)

用途：

- 记录 Go API 的数据结构、接口、验证规则和未来数据库映射。
- 固定前后端契约，避免前端视觉改版时反复改 API。
- 说明 `/api/events`、`/api/cd-items`、`/api/contents`、`/api/connect` 的返回逻辑。
- 说明数据库选型、ORM 方案和第一版表结构。

## 3. 前端多版本设计语言

- [frontend-design-variants.md](frontend-design-variants.md)

用途：

- 记录 V2、V3 以及未来 V4/V5 的设计语言。
- 说明右上角 Review 下拉框如何选择不同版本。
- 把“设计会变”这件事和“业务逻辑不变”分开。

## 4. 本地开发与启动

- [local-development.md](local-development.md)

用途：

- 说明本地启动方式。
- 记录当前可用 API。
- 说明三语言、本地预览和暂缓登录注册。

## 5. 服务器与部署

- [alicloud-tokyo-p0-status.md](alicloud-tokyo-p0-status.md)
- [alicloud-tokyo-p0-deployment.md](alicloud-tokyo-p0-deployment.md)
- [cicd-git-push-deploy.md](cicd-git-push-deploy.md)

用途：

- 记录阿里云东京轻量服务器 P0 配置状态。
- 记录未来部署到 `/opt/livelife` 的路线。
- 和本地开发文档分开，避免把本机调试命令误用到服务器。

## 6. 阶段审批稿

- [product-architecture-and-ui-approval.md](product-architecture-and-ui-approval.md)
- [v3-design-approval.md](v3-design-approval.md)

用途：

- 记录 V2 和 V3 的设计审批过程。
- 不作为唯一开发入口，开发时优先看上面的长期文档。

## 7. 建议新增但待审批的文档

这些文档建议后续补齐，当前先不创建完整内容，等你审批：

- `content-data-entry-guide.md`：后台还没做之前，先规定演出、CD、黑胶、Archive 内容怎么填写。
- `external-link-policy.md`：规定票站、BASE、Instagram、小红书等外部链接怎么放，避免用户误解为站内购买。
- `review-mode-guide.md`：专门说明客户 Review 页面如何切换 V2/V3/V4，哪些控件只在预览环境出现。
- `brand-language-guide.md`：规定 LIVE LIFE、三语言、英文全大写、乐队名纹理和版权/授权提示的用语边界。
- `future-auth-checkout-notes.md`：等决定是否做登录注册和站内支付后，再整理账户、订单、售后的设计。
