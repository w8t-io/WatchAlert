
<p align="center">
  <a href="https://github.com/w8t-io/WatchAlert"> 
    <img src="WatchAlert.png" alt="cloud native monitoring" width="200" height="auto" /></a>
</p>

<p align="center">
  <b>开源一站式监控告警管理系统</b>
</p>

<p align="center">
<a href="https://github.com/w8t-io/WatchAlert/graphs/contributors">
  <img alt="GitHub contributors" src="https://img.shields.io/github/contributors-anon/w8t-io/WatchAlert"/></a>
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/w8t-io/WatchAlert">
<img alt="GitHub forks" src="https://img.shields.io/github/forks/w8t-io/WatchAlert">
<img alt="GitHub Repo issues" src="https://img.shields.io/github/issues/w8t-io/WatchAlert">
<img alt="GitHub Repo issues closed" src="https://img.shields.io/github/issues-closed/w8t-io/WatchAlert">
<img alt="License" src="https://img.shields.io/badge/license-Apache--2.0-blue"/>

## 1. 概述
> WatchAlert作为一款功能强大的开源监控告警管理系统，凭借其丰富的功能、灵活的配置和高效的性能，在云原生和DevOps领域具有广泛的应用前景。无论是对于初创企业还是大型企业，WatchAlert都能提供稳定可靠的监控告警解决方案，助力用户提升运维效率，降低维护成本。
- 开发基础：WatchAlert采用Go语言和React框架开发，确保系统的稳定性和高效性。
- 功能定位：作为一款云原生监控报警管理系统，WatchAlert可以替代AlertManager、PrometheusAlert等组件，提供丰富的监控告警功能。

## 2. 快速开始
- Install Docs：https://cairry.github.io/docs/
- Demo URL：http://8.147.234.89/login
（admin/123）


## 3. 功能介绍
### 人员组织
- 用户管理：包括基本操作和角色绑定，用于管理系统用户。
- 角色管理：包括用户角色的基本操作和通过权限授权实现用户访问控制。

### 告警管理
- 告警规则：支持多数据源和分组通知，用于定义告警规则。
- 静默规则：根据当前告警配置的规则进行告警静默。
- 当前告警：查询当前时间触发的告警列表。
- 历史告警：查询已恢复的历史告警信息。
- 规则模版：内置一些常用的告警规则配置。

### 告警通知
- 通知对象：支持多种通知类型，如飞书（支持官方高级消息卡片Json）、钉钉、企业微信、Email邮件，并配置实际通知模板。
- 通知模板：默认提供3种告警模板，支持创建、更新、删除等基本操作。

### 值班管理
- 值班日程：安排指定成员在特定日期和时间段内处理告警，有效管理告警并提高工作效率。

### 租户管理
- 多租户：允许多个组织共享相同的应用程序实例，同时保持各自数据的隔离性和安全性。

### 数据源
- 提供告警指标的入口，支持多数据源的基本操作。
- 支持Prometheus、阿里云SLS、Loki、Jaeger、AWS CloudWatch作为数据源。

### 仪表盘
- 支持对接 Grafana 面板。

### 日志审计
- 记录重要的操作行为，便于后续审计。

## 交流渠道
- [点击我](https://cairry.github.io/docs/#%E4%BA%A4%E6%B5%81%E7%BE%A4-%E8%81%94%E7%B3%BB%E6%88%91)
