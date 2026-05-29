#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
收银台配置管理平台 - 需求文档生成器
用法: python 生成需求文档.py
"""

import re
import yaml
import sys
from datetime import datetime


def load_template(path="template.yaml"):
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    # 移除 YAML 中的纯注释行（保留行尾注释在解析时处理）
    return yaml.safe_load(content)


def safe_get(data, *keys, default=""):
    """安全地从嵌套字典中获取值"""
    for key in keys:
        if isinstance(data, dict):
            data = data.get(key, {})
        else:
            return default
    return data if data is not None else default


def gen_table(headers, rows):
    """生成 Markdown 表格"""
    if not rows:
        return ""
    sep = "|" + "|".join("---" for _ in headers) + "|"
    header = "| " + " | ".join(headers) + " |"
    lines = [header, sep]
    for row in rows:
        lines.append("| " + " | ".join(str(c) for c in row) + " |")
    return "\n".join(lines)


def priority_label(p):
    labels = {"P0": "🔴 核心（P0）", "P1": "🟡 重要（P1）", "P2": "🟢 优化（P2）"}
    return labels.get(p, p)


def main():
    print("正在解析模板...")
    data = load_template()

    proj = data.get("project", {})
    biz = data.get("business", {})
    modules = data.get("modules", {})
    roles = data.get("roles", [])
    routes = data.get("routes", [])
    db = data.get("database", {})
    ts = data.get("tech_stack", {})
    nf = data.get("non_functional", {})
    notes = data.get("notes", [])

    doc = []
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # ==================== 标题 ====================
    doc.append(f"# {proj.get('name', '收银台配置管理平台')} v{proj.get('version', '1.0.0')}")
    doc.append("")
    doc.append(f"> 本文档由需求模板自动生成，生成时间：{now}")
    doc.append("")
    doc.append("---")
    doc.append("")

    # ==================== 1. 项目概述 ====================
    doc.append("## 1. 项目概述")
    doc.append("")

    overview_rows = [
        ("项目名称", proj.get("name", "")),
        ("版本号", proj.get("version", "")),
        ("负责部门", proj.get("department", "")),
        ("项目描述", proj.get("description", "")),
    ]
    doc.append(gen_table(("项目", "内容"), overview_rows))
    doc.append("")

    bg = biz.get("background", "")
    if bg:
        doc.append("### 1.1 项目背景")
        doc.append("")
        doc.append(bg.strip())
        doc.append("")

    goals = biz.get("goals", [])
    if goals:
        doc.append("### 1.2 项目目标")
        doc.append("")
        for g in goals:
            doc.append(f"- {g}")
        doc.append("")

    users = biz.get("target_users", [])
    if users:
        doc.append("### 1.3 目标用户")
        doc.append("")
        for u in users:
            doc.append(f"- {u}")
        doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 2. 用户角色与权限 ====================
    doc.append("## 2. 用户角色与权限")
    doc.append("")

    role_rows = []
    for r in roles:
        perms = r.get("permissions", [])
        role_rows.append((
            r.get("name", ""),
            r.get("desc", ""),
            "、".join(perms) if isinstance(perms, list) else str(perms),
        ))
    doc.append(gen_table(("角色", "描述", "权限"), role_rows))
    doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 3. 前端功能模块 ====================
    doc.append("## 3. 前端功能模块")
    doc.append("")

    frontend = modules.get("frontend", [])
    if isinstance(frontend, list):
        for mod in frontend:
            p = priority_label(mod.get("priority", ""))
            doc.append(f"### {mod.get('name', '')} — {p}")
            doc.append("")
            doc.append(f"{mod.get('description', '')}")
            doc.append("")
            features = mod.get("features", [])
            if features:
                doc.append("**功能列表：**")
                doc.append("")
                for f in features:
                    doc.append(f"- {f}")
                doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 4. 后端 API 模块 ====================
    doc.append("## 4. 后端 API 接口")
    doc.append("")

    backend = modules.get("backend", [])
    if isinstance(backend, list):
        for mod in backend:
            p = priority_label(mod.get("priority", ""))
            doc.append(f"### {mod.get('name', '')} — {p}")
            doc.append("")
            doc.append(f"{mod.get('description', '')}")
            doc.append("")
            endpoints = mod.get("endpoints", [])
            if endpoints:
                doc.append("**API 列表：**")
                doc.append("")
                ep_rows = []
                for ep in endpoints:
                    if isinstance(ep, str):
                        m = re.match(r"(GET|POST|PUT|DELETE|PATCH)\s+(\S+)\s*[#，]\s*(.*)", ep)
                        if m:
                            ep_rows.append((m.group(1), m.group(2), m.group(3)))
                        else:
                            ep_rows.append(("", ep, ""))
                    elif isinstance(ep, dict):
                        ep_rows.append((
                            ep.get("method", ""),
                            ep.get("path", ""),
                            ep.get("desc", ""),
                        ))
                if ep_rows:
                    doc.append(gen_table(("方法", "路径", "说明"), ep_rows))
                doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 5. 数据库设计 ====================
    doc.append("## 5. 数据库设计")
    doc.append("")

    tables = db.get("tables", [])
    if tables:
        tbl_rows = []
        for t in tables:
            tbl_rows.append((
                t.get("name", ""),
                t.get("comment", ""),
                t.get("fields", ""),
            ))
        doc.append(gen_table(("表名", "说明", "主要字段"), tbl_rows))
        doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 6. 技术栈 ====================
    doc.append("## 6. 技术栈")
    doc.append("")

    tech_labels = {
        "frontend": {
            "framework": "框架",
            "ui_library": "UI 组件库",
            "build_tool": "构建工具",
            "state_management": "状态管理",
            "http_client": "HTTP 客户端",
            "csv_library": "CSV 处理库",
            "form_library": "表单库",
            "additional": "其他",
        },
        "backend": {
            "language": "开发语言",
            "framework": "HTTP 框架",
            "router": "路由",
            "orm": "ORM",
            "database": "数据库",
            "cache": "缓存",
            "migration": "数据库迁移",
            "config": "配置管理",
            "auth": "认证与权限",
            "csv": "CSV 处理",
            "validator": "参数校验",
            "api_docs": "API 文档",
            "log": "日志库",
            "message_queue": "消息队列",
            "test": "测试框架",
        },
        "deployment": {
            "container": "容器化",
            "ci_cd": "CI/CD",
            "server": "服务器",
            "reverse_proxy": "反向代理",
        },
    }

    for section in ["frontend", "backend", "deployment"]:
        section_data = ts.get(section, {})
        if section_data:
            section_name = {"frontend": "前端", "backend": "后端", "deployment": "部署运维"}.get(section, section)
            idx = {"前端": "1", "后端": "2", "部署运维": "3"}.get(section_name, "4")
            doc.append(f"### 6.{idx} {section_name}技术栈")
            doc.append("")
            rows = []
            for k, v in section_data.items():
                label = tech_labels.get(section, {}).get(k, k)
                if v:
                    rows.append((label, str(v)))
            if rows:
                doc.append(gen_table(("类别", "选型"), rows))
                doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 7. 非功能性需求 ====================
    doc.append("## 7. 非功能性需求")
    doc.append("")

    nf_labels = {
        "performance": "性能要求",
        "security": "安全要求",
        "availability": "可用性要求",
        "scalability": "可扩展性要求",
    }

    for key, label in nf_labels.items():
        items = nf.get(key, [])
        if items:
            doc.append(f"### {label}")
            doc.append("")
            for item in items:
                doc.append(f"- {item}")
            doc.append("")

    doc.append("---")
    doc.append("")

    # ==================== 8. 页面路由 ====================
    doc.append("## 8. 页面路由规划")
    doc.append("")

    if routes:
        route_rows = []
        for r in routes:
            role_list = r.get("roles", [])
            if isinstance(role_list, list):
                role_str = "、".join(role_list)
            else:
                role_str = str(role_list)
            route_rows.append((
                r.get("path", ""),
                r.get("name", ""),
                r.get("component", ""),
                role_str,
            ))
        doc.append(gen_table(("路径", "页面名称", "组件", "可访问角色"), route_rows))
        doc.append("")

    # ==================== 9. 附加说明 ====================
    if notes:
        doc.append("---")
        doc.append("")
        doc.append("## 9. 附加说明")
        doc.append("")
        for n in notes:
            doc.append(f"> - {n}")

    doc.append("")
    doc.append("---")
    doc.append("")
    doc.append(f"*文档结束 — 生成时间：{now}*")
    doc.append("")

    output = "\n".join(doc)

    with open("需求文档.md", "w", encoding="utf-8") as f:
        f.write(output)

    print(f"[OK] 需求文档已生成: 需求文档.md")


if __name__ == "__main__":
    main()
