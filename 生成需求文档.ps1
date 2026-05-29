# 收银台配置管理平台 - 需求文档生成器
# 用法: pwsh .\生成需求文档.ps1

param(
    [string]$TemplatePath = "template.yaml",
    [string]$OutputPath = "需求文档.md"
)

# 检查模板文件
if (-not (Test-Path $TemplatePath)) {
    Write-Error "未找到模板文件: $TemplatePath"
    exit 1
}

# 读取模板 YAML
$yaml = Get-Content $TemplatePath -Raw

# 简单 YAML 解析函数（处理嵌套结构）
function Parse-Yaml {
    param([string]$Content)
    
    $lines = $Content -split "`n"
    $result = @{}
    $stack = @()
    $currentIndent = 0
    
    foreach ($line in $lines) {
        if ($line -match "^\s*#" -or $line.Trim() -eq "") { continue }
        
        $indent = ($line | Select-String "^(\s*)").Matches[0].Groups[1].Value.Length
        
        # 移除行内注释
        $lineClean = $line -replace "\s+#.*$", ""
        
        if ($lineClean -match "^\s*(\w[\w\d_-]*)\s*:\s*(.*)") {
            $key = $matches[1]
            $value = $matches[2].Trim()
            
            # 处理缩进变化
            while ($stack.Count -gt 0 -and $indent -le $stack[-1].indent) {
                $null = $stack.Pop()
            }
            
            if ($value -eq "" -or $value -eq ">") {
                # 对象或数组开始
                $entry = @{ key = $key; indent = $indent; children = @{} }
                $stack += $entry
            }
            elseif ($value -match "^\s*-\s*(.*)") {
                # 列表项
                $itemValue = $matches[1].Trim()
            }
            else {
                # 简单值
                if ($stack.Count -gt 0) {
                    $parent = $stack[-1]
                    if (-not $parent.children.ContainsKey($key)) {
                        $parent.children[$key] = $value
                    }
                }
                else {
                    $result[$key] = $value
                }
            }
        }
        elseif ($lineClean -match "^\s*-\s+(.*)") {
            $itemValue = $matches[1].Trim()
            # 列表项 - 简化处理
        }
    }
    
    # 简化：使用 ConvertFrom-Yaml (PowerShell 7+)
    return $result
}

# 使用 PowerShell 的 Yaml 处理（需要 PowershellGet）
# 尝试用 ConvertFrom-Json 配合简单替换处理
function Convert-YamlToObject {
    param([string]$Content)
    
    # 首先尝试用 YamlDotNet
    try {
        $assembly = [System.Reflection.Assembly]::LoadWithPartialName("YamlDotNet") 2>$null
        if (-not $assembly) {
            # 尝试加载 NuGet 包
            $null = $null
        }
    }
    catch {}
    
    # 用简化的方案：将 YAML 用分隔符包围，分段处理
    # 读取并结构化为嵌套哈希表
    
    $result = @{}
    $currentPath = @()
    $indentStack = @()
    $currentList = $null
    $currentListItem = $null
    $listItemStack = @()
    $inMultilineValue = $false
    $multilineKey = ""
    $multilineLines = @()
    
    $lines = $Content -split "`n"
    
    foreach ($rawLine in $lines) {
        $line = $rawLine
        
        # 跳过注释和空行
        if ($line.TrimStart() -match "^#" -or $line.Trim() -eq "") {
            if ($inMultilineValue) {
                $multilineLines += ""
            }
            continue
        }
        
        # 计算缩进
        $indent = if ($line -match "^(\s*)") { $matches[1].Length } else { 0 }
        
        # 处理多行值 (> 语法)
        if ($inMultilineValue) {
            if ($indent -ge 4) {
                $multilineLines += $line.Trim()
                continue
            }
            else {
                # 多行值结束
                $inMultilineValue = $false
                $value = $multilineLines -join " "
                Set-YamlValue $result $currentPath $multilineKey $value
                $multilineLines = @()
            }
        }
        
        $trimmedLine = $line.Trim()
        
        # 跳过纯列表行（子项会被其 children 处理）
        if ($trimmedLine -match "^\s*-\s+\w") {
            # 列表项，尝试提取 key: value
            if ($trimmedLine -match "^\s*-\s+([^:]+):\s*(.*)") {
                $value = $matches[2].Trim()
                if ($value -eq "") {
                    # 列表项包含子属性
                }
            }
            continue
        }
        
        # 匹配 key: value 或 key:
        if ($trimmedLine -match "^(\w[\w\d_-]*)\s*:\s*(.*)") {
            $key = $matches[1]
            $value = $matches[2].Trim()
            
            # 更新缩进路径
            while ($indentStack.Count -gt 0 -and $indent -le $indentStack[-1].indent) {
                $null = $indentStack.Pop()
                if ($currentPath.Count -gt 0) { $currentPath = $currentPath[0..($currentPath.Count-2)] }
            }
            
            $indentStack += @{ indent = $indent; key = $key }
            $currentPath += $key
            
            if ($value -eq ">" -or $value -eq "") {
                # 对象，继续
            }
            else {
                Set-YamlValue $result $currentPath $key $value
            }
        }
    }
    
    return $result
}

function Set-YamlValue {
    param($hash, $path, $key, $value)
    
    if ($path.Count -eq 0) { return }
    
    $current = $hash
    for ($i = 0; $i -lt $path.Count - 1; $i++) {
        if (-not $current.ContainsKey($path[$i])) {
            $current[$path[$i]] = @{}
        }
        $current = $current[$path[$i]]
    }
    
    $parentKey = $path[-1]
    if (-not $current.ContainsKey($parentKey)) {
        $current[$parentKey] = @{}
    }
    $current = $current[$parentKey]
    
    if ($key -ne $parentKey) {
        $current[$key] = $value
    }
    else {
        # 已经是同一个 key
    }
}

# 由于纯 PowerShell YAML 解析较复杂，这里提供一个内嵌的解析方案
# 改用更可靠的方式：逐段手动解析我们已知的结构

# 实际使用更稳健的方案 - 直接用对象表示法
$templateContent = Get-Content $TemplatePath -Raw

# 从 YAML 中提取项目信息（使用正则逐段提取）
function Get-YamlValueByPath {
    param([string]$Content, [string]$Path)
    
    $keys = $Path -split '\.'
    $pattern = ''
    foreach ($key in $keys) {
        if ($pattern -eq '') {
            $pattern = "^$key\s*:\s*(.*)"
        } else {
            $pattern = "^\s+$key\s*:\s*(.*)"
        }
    }
    
    if ($Content -match $pattern) {
        return $matches[1].Trim()
    }
    return $null
}

# ================= 开始生成文档 =================

Write-Host "正在生成需求文档..." -ForegroundColor Cyan

# 提取关键信息
$projectName = "收银台配置管理平台"
$projectVersion = "1.0.0"
$projectDesc = ""
$projectDept = ""

if ($templateContent -match "^name:\s*""(.+)""") { $projectName = $matches[1] }
elseif ($templateContent -match "^name:\s*(.+)") { $projectName = $matches[1].Trim('" ') }
if ($templateContent -match "^version:\s*""(.+)""") { $projectVersion = $matches[1] }
if ($templateContent -match "^description:\s*""(.+)""") { $projectDesc = $matches[1] }
if ($templateContent -match "^department:\s*""(.+)""") { $projectDept = $matches[1] }

$background = ""
if ($templateContent -match "^background:\s*>\s*(.+?)(?=^\S|\Z)") { $background = $matches[1] -replace "`n", " " -replace "\s+", " " }

# 提取目标
$goals = @()
$inGoals = $false
foreach ($line in ($templateContent -split "`n")) {
    if ($line -match "^\s+goals:") { $inGoals = $true; continue }
    if ($inGoals -and $line -match "^\s+- ") { $goals += $line -replace "^\s+-\s+""?", "" -replace """$", "" }
    if ($inGoals -and $line -match "^\s+\w" -and $line -notmatch "^\s+- ") { break }
}

$targetUsers = @()
$inUsers = $false
foreach ($line in ($templateContent -split "`n")) {
    if ($line -match "^\s+target_users:") { $inUsers = $true; continue }
    if ($inUsers -and $line -match "^\s+- ") { $targetUsers += $line -replace "^\s+-\s+""?", "" -replace """$", "" }
    if ($inUsers -and $line -match "^\s+\w" -and $line -notmatch "^\s+- ") { break }
}

# 解析模块
$frontendModules = @()
$backendModules = @()
$roles = @()
$routes = @()
$tables = @()
$techStack = @{}
$nonFunc = @()

$currentSection = $null
$currentModule = $null
$currentFeatures = $null
$currentEndpoints = $null
$inPermissions = $false
$tableName = ""
$tableComment = ""
$tableFields = ""

foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    $indent = if ($line -match "^(\s*)") { $matches[1].Length } else { 0 }
    
    # 跳过注释
    if ($trimmed -match "^#") { continue }
    if ($trimmed -eq "") { continue }
    
    # 检测顶级 section
    if ($trimmed -match "^(frontend|backend):" -and $indent -eq 2) {
        $currentSection = $matches[1]
        continue
    }
    
    # 在 frontend/backend 下检测模块
    if (($currentSection -eq "frontend" -or $currentSection -eq "backend") -and $indent -eq 6) {
        if ($trimmed -match "^- name:\s*""(.+)""") {
            $moduleName = $matches[1]
        }
        elseif ($trimmed -match "^- name:\s*(.+)") {
            $moduleName = $matches[1].Trim(' "')
        }
        else { continue }
        
        $currentModule = @{
            name = $moduleName
            desc = ""
            priority = ""
            features = @()
            endpoints = @()
        }
        
        if ($currentSection -eq "frontend") {
            $frontendModules += $currentModule
        }
        else {
            $backendModules += $currentModule
        }
        continue
    }
    
    if ($currentModule) {
        if ($trimmed -match "^description:\s*""(.+)""") {
            $currentModule.desc = $matches[1]
        }
        elseif ($trimmed -match "^priority:\s*""(.+)""") {
            $currentModule.priority = $matches[1]
        }
        elseif ($trimmed -match "^\s+-\s+""(.+)""") {
            $feat = $matches[1]
            if ($currentModule.features -is [array]) { $currentModule.features += $feat }
        }
        elseif ($trimmed -match "^\s+-\s+(.+?)\s+#") {
            $ep = $matches[1].Trim()
            if ($currentModule.endpoints -is [array]) { $currentModule.endpoints += $ep }
        }
        elseif ($trimmed -match "^\s+-\s+(.+)") {
            $val = $matches[1].Trim()
            if ($val -match "^(GET|POST|PUT|DELETE|PATCH)") {
                if ($currentModule.endpoints -is [array]) { $currentModule.endpoints += $val }
            }
            else {
                if ($currentModule.features -is [array]) { $currentModule.features += $val }
            }
        }
    }
    
    # 解析 roles (没有可靠的 YAML 解析时的替代方案)
    if ($trimmed -eq "roles:" -and $indent -eq 0) {
        $currentSection = "roles"
        continue
    }
}

# 重新读取 roles 段落
$inRolesSection = $false
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    $indent = if ($line -match "^(\s*)") { $matches[1].Length } else { 0 }
    
    if ($trimmed -eq "roles:" -and $indent -eq 0) {
        $inRolesSection = $true
        continue
    }
    if ($inRolesSection -and $trimmed -match "^modules:") { $inRolesSection = $false; continue }
    
    if ($inRolesSection -and $trimmed -match "^- name:\s*""(.+)""") {
        $roles += @{ name = $matches[1]; desc = ""; permissions = @() }
    }
    elseif ($inRolesSection -and $trimmed -match "^- name:\s*(.+)") {
        $roles += @{ name = $matches[1].Trim(' "'); desc = ""; permissions = @() }
    }
    elseif ($inRolesSection -and $trimmed -match "^desc:\s*""(.+)""") -and $roles.Count -gt 0) {
        $roles[-1].desc = $matches[1]
    }
    elseif ($inRolesSection -and $trimmed -match "^desc:\s*(.+)") -and $roles.Count -gt 0) {
        $roles[-1].desc = $matches[1].Trim(' "')
    }
    elseif ($inRolesSection -and $trimmed -match "^\s+-\s+""(.+)""") -and $roles.Count -gt 0) {
        $roles[-1].permissions += $matches[1]
    }
    elseif ($inRolesSection -and $trimmed -match "^\s+-\s+(.+)") -and $roles.Count -gt 0) {
        $roles[-1].permissions += $matches[1].Trim(' "')
    }
}

# 重新读取 routes
$inRoutes = $false
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    
    if ($trimmed -eq "routes:" -and !$inRoutes) {
        $inRoutes = $true
        continue
    }
    if ($inRoutes -and $trimmed -match "^- path:") {
        $route = @{ path = ""; name = ""; component = ""; roles = @() }
        if ($trimmed -match 'path:\s*"(.+)"') { $route.path = $matches[1] }
        $routes += $route
    }
    elseif ($inRoutes -and $trimmed -match '^name:\s*"(.+)"' -and $routes.Count -gt 0) {
        $routes[-1].name = $matches[1]
    }
    elseif ($inRoutes -and $trimmed -match '^component:\s*"(.+)"' -and $routes.Count -gt 0) {
        $routes[-1].component = $matches[1]
    }
    elseif ($inRoutes -and $trimmed -match '^roles:\s*\[(.+)\]' -and $routes.Count -gt 0) {
        $routes[-1].roles = $matches[1] -split ",\s*" | ForEach-Object { $_.Trim(' "') }
    }
    elseif ($inRoutes -and $trimmed -match "^notes:") { break }
}

# 读取数据库表
$inTables = $false
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    
    if ($trimmed -eq "tables:") { $inTables = $true; continue }
    if ($inTables -and $trimmed -match "^- name:") {
        if ($tableName -ne "") { $tables += @{ name = $tableName; comment = $tableComment; fields = $tableFields } }
        $tableName = if ($trimmed -match 'name:\s*"(.+)"') { $matches[1] } else { "" }
        $tableComment = ""; $tableFields = ""
    }
    elseif ($inTables -and $trimmed -match '^comment:\s*"(.+)"') { $tableComment = $matches[1] }
    elseif ($inTables -and $trimmed -match '^fields:\s*"(.+)"') { $tableFields = $matches[1] }
    elseif ($inTables -and $trimmed -match "^notes:" -or $trimmed -match "^tech_stack:") { $inTables = $false }
}
if ($tableName -ne "") { $tables += @{ name = $tableName; comment = $tableComment; fields = $tableFields } }

# 读取技术栈
$inTech = $false; $techCategory = ""
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    $indent = if ($line -match "^(\s*)") { $matches[1].Length } else { 0 }
    
    if ($trimmed -eq "tech_stack:") { $inTech = $true; continue }
    if ($inTech -and $trimmed -match "^(frontend|backend|deployment):") { $techCategory = $matches[1]; continue }
    if ($inTech -and $trimmed -match "^non_functional:") { break }
    
    if ($inTech -and $techCategory -and $trimmed -match "^(\w[\w\d_-]*):\s*(.*)") {
        $key = $matches[1]; $val = $matches[2].Trim(' "')
        if (-not $techStack[$techCategory]) { $techStack[$techCategory] = @{} }
        $techStack[$techCategory][$key] = $val
    }
}

# 读取非功能需求
$currentNF = ""
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    $indent = if ($line -match "^(\s*)") { $matches[1].Length } else { 0 }
    
    if ($trimmed -eq "non_functional:") { $currentNF = "nf"; continue }
    if ($currentNF -eq "nf" -and $indent -eq 2 -and $trimmed -match "^(\w[\w\d_-]*):") {
        $currentNF = $matches[1]; $nonFunc[$currentNF] = @(); continue
    }
    if ($currentNF -eq "nf" -and $trimmed -match "^routes:") { break }
    
    if ($currentNF -ne "nf" -and $currentNF -ne "" -and $trimmed -match "^\s+-\s+") {
        $val = $trimmed -replace "^\s+-\s+", "" -replace '"', ""
        $nonFunc[$currentNF] += $val
    }
}

Write-Host "数据解析完成，开始生成文档..." -ForegroundColor Green

# ================= 生成 Markdown 文档 =================

$doc = @"

# $projectName v$projectVersion

> 本文档由需求模板自动生成，生成时间：$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

---

## 1. 项目概述

| 项目 | 内容 |
|------|------|
| 项目名称 | $projectName |
| 版本号 | $projectVersion |
| 负责部门 | $projectDept |
| 项目描述 | $projectDesc |

### 1.1 项目背景

$background

### 1.2 项目目标

"@

foreach ($goal in $goals) {
    $doc += @"
- $goal
"@
}

$doc += @"

### 1.3 目标用户

"@

foreach ($user in $targetUsers) {
    $doc += @"
- $user
"@
}

$doc += @"

---

## 2. 用户角色与权限

| 角色 | 描述 | 权限 |
|------|------|------|
"@

foreach ($role in $roles) {
    $perms = $role.permissions -join "、"
    $doc += @"
| $($role.name) | $($role.desc) | $perms |
"@
}

$doc += @"

---

## 3. 前端功能模块

"@

foreach ($mod in $frontendModules) {
    $priorityLabel = switch ($mod.priority) {
        "P0" { "🔴 核心（P0）" }
        "P1" { "🟡 重要（P1）" }
        "P2" { "🟢 优化（P2）" }
        default { $mod.priority }
    }
    $doc += @"
### $($mod.name) — $priorityLabel

$($mod.desc)

**功能列表：**
"@
    foreach ($feat in $mod.features) {
        $doc += @"
- $feat
"@
    }
    $doc += @"

"@
}

$doc += @"
---

## 4. 后端 API 模块

"@

foreach ($mod in $backendModules) {
    $priorityLabel = switch ($mod.priority) {
        "P0" { "🔴 核心（P0）" }
        "P1" { "🟡 重要（P1）" }
        "P2" { "🟢 优化（P2）" }
        default { $mod.priority }
    }
    $doc += @"
### $($mod.name) — $priorityLabel

$($mod.desc)

**API 列表：**

| 方法 | 路径 | 说明 |
|------|------|------|
"@
    foreach ($ep in $mod.endpoints) {
        if ($ep -match "^(GET|POST|PUT|DELETE|PATCH)\s+(\S+)\s+#\s*(.*)") {
            $doc += @"
| $($matches[1]) | $($matches[2]) | $($matches[3]) |
"@
        }
        elseif ($ep -match "^(GET|POST|PUT|DELETE|PATCH)\s+(\S+)\s*#\s*(.*)") {
            $doc += @"
| $($matches[1]) | $($matches[2]) | $($matches[3]) |
"@
        }
    }
    $doc += @"

"@
}

$doc += @"
---

## 5. 数据库设计

| 表名 | 说明 | 主要字段 |
|------|------|----------|
"@

foreach ($table in $tables) {
    $doc += @"
| $($table.name) | $($table.comment) | $($table.fields) |
"@
}

$doc += @"

---

## 6. 技术栈

### 6.1 前端技术栈

| 类别 | 选型 |
|------|------|
"@

if ($techStack.ContainsKey("frontend")) {
    foreach ($kv in $techStack["frontend"].GetEnumerator()) {
        $label = switch ($kv.Key) {
            "framework" { "框架" }
            "ui_library" { "UI 组件库" }
            "build_tool" { "构建工具" }
            "state_management" { "状态管理" }
            "http_client" { "HTTP 客户端" }
            "additional" { "其他" }
            default { $kv.Key }
        }
        $doc += @"
| $label | $($kv.Value) |
"@
    }
}

$doc += @"

### 6.2 后端技术栈

| 类别 | 选型 |
|------|------|
"@

if ($techStack.ContainsKey("backend")) {
    foreach ($kv in $techStack["backend"].GetEnumerator()) {
        $label = switch ($kv.Key) {
            "language" { "开发语言" }
            "framework" { "框架" }
            "orm" { "ORM" }
            "database" { "数据库" }
            "cache" { "缓存" }
            "message_queue" { "消息队列" }
            "api_docs" { "API 文档" }
            default { $kv.Key }
        }
        $doc += @"
| $label | $($kv.Value) |
"@
    }
}

$doc += @"

### 6.3 部署运维

| 类别 | 选型 |
|------|------|
"@

if ($techStack.ContainsKey("deployment")) {
    foreach ($kv in $techStack["deployment"].GetEnumerator()) {
        $label = switch ($kv.Key) {
            "container" { "容器化" }
            "ci_cd" { "CI/CD" }
            "server" { "服务器" }
            default { $kv.Key }
        }
        $doc += @"
| $label | $($kv.Value) |
"@
    }
}

$doc += @"

---

## 7. 非功能性需求

"@

foreach ($nf in $nonFunc.GetEnumerator()) {
    $nfLabel = switch ($nf.Key) {
        "performance" { "性能要求" }
        "security" { "安全要求" }
        "availability" { "可用性要求" }
        "scalability" { "可扩展性要求" }
        default { $nf.Key }
    }
    $doc += @"
### $nfLabel

"@
    foreach ($item in $nf.Value) {
        $doc += @"
- $item
"@
    }
    $doc += @"

"@
}

$doc += @"
---

## 8. 页面路由规划

| 路径 | 页面名称 | 组件 | 可访问角色 |
|------|----------|------|-----------|
"@

foreach ($route in $routes) {
    $roleStr = $route.roles -join "、"
    $doc += @"
| $($route.path) | $($route.name) | $($route.component) | $roleStr |
"@
}

# 查找附加说明
$notes = @()
$inNotes = $false
foreach ($line in ($templateContent -split "`n")) {
    $trimmed = $line.Trim()
    if ($trimmed -eq "notes:") { $inNotes = $true; continue }
    if ($inNotes -and $trimmed -match "^- ") {
        $notes += $trimmed -replace "^- ", "" -replace '"', ""
    }
    if ($inNotes -and $trimmed -notmatch "^- " -and $trimmed -ne "") { break }
}

if ($notes.Count -gt 0) {
    $doc += @"

---

## 9. 附加说明

"@
    foreach ($note in $notes) {
        $doc += @">
- $note
"@
    }
}

$doc += @"

---

*文档结束 — 生成时间：$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')*
"@

# 写入文件
$doc | Set-Content $OutputPath -Encoding UTF8

Write-Host "✅ 需求文档已生成: $OutputPath" -ForegroundColor Green
Write-Host "  文件路径: $(Resolve-Path $OutputPath)" -ForegroundColor Green
