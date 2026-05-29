#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
CSV 文件生成器
从 YAML 数据源生成与样例完全一致的 CSV 文件

用法:
  python csv_generator.py                          # 使用内置样例数据
  python csv_generator.py -i data.yaml             # 使用自定义 YAML 数据
  python csv_generator.py -i data.yaml -o ./output # 指定输出目录
"""

import csv
import json
import os
import sys
import argparse
from datetime import datetime

try:
    import yaml
except ImportError:
    print("请先安装 PyYAML: pip install pyyaml")
    sys.exit(1)

# ==================== 样例数据（完全匹配 csv/ 下的格式） ====================

SAMPLE_MODULES = [
    ["AllInOne_JPOS", "CN", '[{"name": "AllInOne_JPOS", "steps": [{"name": "HWM", "path": "HWM/CN/ALLINONE/hwm_jpos.ps1"}, {"name": "UPOS", "path": "RFID/CN/UPOS.ps1"}, {"name": "JPOS", "path": "JavaPos/CN/ALLINONE/jpos.ps1"}, {"name": "EFT", "path": "EFT/CN/ICBC/icbc.ps1"}]}, {"name": "Cashtill_Toolbox", "steps": [{"name": "hotkey", "path": "Toolbox/HotKey/hot_key.ps1"}, {"name": "AsiaCashtillTool", "path": "Toolbox/AsiaWindowsCashtillTool/cashtill_tool.ps1"}]}]', "pr", 7],
    ["M2", "CN", '[{"name": "M2", "steps": [{"name": "HWM", "path": "HWM/CN/Till/hwm.ps1"}, {"name": "UPOS", "path": "RFID/CN/UPOS.ps1"}, {"name": "JPOS", "path": "JavaPos/CN/M2&R2/th200i.ps1"}, {"name": "EFT", "path": "EFT/CN/ICBC/icbc.ps1"}]}, {"name": "Cashtill_Toolbox", "steps": [{"name": "hotkey", "path": "Toolbox/HotKey/hot_key.ps1"}, {"name": "AsiaCashtillTool", "path": "Toolbox/AsiaWindowsCashtillTool/cashtill_tool.ps1"}]}]', "pr", 8],
    ["M2_TH200i", "PH", '[{"name": "M2_TH200i", "steps": [{"name": "hwm", "path": "HWM/M2_TH200i_Seazone.ps1"}, {"name": "javapos", "path": "JavaPos/th200i.ps1"}, {"name": "upos", "path": "RFID/UPOS.ps1"}]}, {"name": "PH_EFT", "steps": [{"name": "eft", "path": "EFT/PH/eft.ps1"}]}, {"name": "Cashtill_Toolbox", "steps": [{"name": "hotkey", "path": "Toolbox/HotKey/hot_key.ps1"}, {"name": "AsiaCashtillTool", "path": "Toolbox/AsiaWindowsCashtillTool/cashtill_tool.ps1"}]}]', "pr", 10],
    ["M2_TH200i", "TH", '[{"name": "M2_TH200i", "steps": [{"name": "hwm", "path": "HWM/M2_TH200i_Seazone.ps1"}, {"name": "javapos", "path": "JavaPos/th200i.ps1"}, {"name": "upos", "path": "RFID/UPOS.ps1"}]}, {"name": "TH_EFT", "steps": [{"name": "eft", "path": "EFT/TH/eft.ps1"}]}, {"name": "Cashtill_Toolbox", "steps": [{"name": "hotkey", "path": "Toolbox/HotKey/hot_key.ps1"}, {"name": "AsiaCashtillTool", "path": "Toolbox/AsiaWindowsCashtillTool/cashtill_tool.ps1"}]}]', "pr", 13],
]

SAMPLE_RULES = [
    ["CN_LAB_Register", "env", "CN", "InfraAuto_Hostname", "cnt0806.*", "lab"],
    ["CN_PP_Store", "env", "CN", "InfraAuto_Hostname", "CNT1628.*|CNT0901.*|CNT0422.*|CNT1778.*|CNT1260.*|CNT1277.*|CNT0861.*|CNT0921.*|CNT1226.*|CNT0860.*|CNT2966.*|CNT0990.*|CNT1265.*|CNT0488.*|CNT1758.*", "pp"],
    ["CN_PP_AllInOne", "env", "CN", "InfraAuto_Hostname", "CNT127749|CNT127708|CNT038806|CNT325802|CNT038805|CNT176722", "pp"],
    ["CN_PR_Register", "env", "CN", "InfraAuto_Hostname", ".*", "pr"],
    ["MY_PR_Register", "env", "MY", "InfraAuto_Hostname", ".*", "pr"],
    ["SG_PR_Register", "env", "SG", "InfraAuto_Hostname", ".*", "pr"],
    ["CN_PR_M2", "group", "CN", "InfraAuto_Device_Model", "BeetleA", "BeetleA_P1200"],
    ["CN_PR_EngagePOS", "group", "CN", "InfraAuto_Device_Model", "EngagePOS", "EngagePOS"],
    ["SG_PR_BeetleA_P1200", "group", "SG", "InfraAuto_Device_Model", "BeetleA", "BeetleA_P1200"],
    ["TW_PR_M2_TH200i", "group", "TW", "InfraAuto_Hostname", ".*", "M2_TH200i"],
]

SAMPLE_STORES = [
    [806, "10.66.48.0", "pr", "ICBC", "CN", "10.69.127.244", "10.69.63.225"],
    [388, "10.5.253.0", "pr", "GMC", "CN", "10.69.96.4", "10.69.0.1"],
    [391, "10.5.249.0", "pr", "GMC", "CN", "10.69.96.36", "10.69.0.65"],
    [392, "10.5.251.0", "pr", "ICBC", "CN", "10.69.96.52", "10.69.0.97"],
    [969, "10.17.253.0", "pr", "NETS", "SG", "10.17.52.212", "10.17.56.209"],
    [2128, "10.17.249.0", "pr", "NETS", "SG", "10.17.53.36", "10.17.57.33"],
    [666, "10.0.253.0", "pr", "Taishin", "TW", "10.0.53.244", "10.0.59.225"],
    [1035, "10.10.5.0", "pr", "None", "TH", "10.10.56.4", "10.10.52.1"],
    [1708, "10.13.253.0", "pr", "BCA", "ID", "10.13.10.20", "10.13.12.17"],
    [1789, "10.18.252.0", "pr", "QUEST", "AU", "10.18.10.4", "10.18.12.1"],
]

SAMPLE_TILLS = [
    ["CNT080634", "4C-52-62-35-BC-01"],
    ["CNT162801", "4C-52-62-35-BC-02"],
    ["CNT090101", "4C-52-62-35-BC-03"],
    ["SGT212801", "00-1E-67-62-35-01"],
    ["TWT066601", "00-1E-67-62-35-02"],
    ["IDT170804", "4C-52-62-35-BC-10"],
]

SAMPLE_VARS = [
    ["TW_PARK_STORE_NO", "1222,1229,1292,1633,1639,1681,1762,1766,1769,1772,1779,666,871,928,988,926", "pr", "[]"],
    ["Greengate_Store", "806", "lab", "[]"],
    ["Greengate_Store", "1277,856,665,1719,1203,392,1265,1670,422,947,1204,419,842,1289", "pr", "[]"],
    ["Greengate_Store", "1277,856,665,1719,1203,392,1265,1670,422,947,1204,419,842,1289", "pp", "[]"],
    ["AsiaWindowsCashtillTool_SRC", "https://infra-auto.dktapp.cloud/public/shell/AsiaWindowsCashtillTool_1.0.14.exe", "lab", "[]"],
    ["AsiaWindowsCashtillTool_SRC", "https://infra-auto.dktapp.cloud/public/shell/AsiaWindowsCashtillTool_1.0.14.exe", "pp", "[]"],
    ["AsiaWindowsCashtillTool_SRC", "https://infra-auto.dktapp.cloud/public/shell/AsiaWindowsCashtillTool_1.0.12.exe", "pr", "[]"],
    ["AsiaWindowsCashtillTool_MD5", "45041EC8238602BE86D31F7E028318D6", "lab", "[]"],
    ["AsiaWindowsCashtillTool_MD5", "45041EC8238602BE86D31F7E028318D6", "pp", "[]"],
    ["AsiaWindowsCashtillTool_MD5", "8084309A56A5B2EB1BC932D34C53F844", "pr", "[]"],
    ["AsiaWindowsCashtillTool_Ver", "1.0.14", "lab", "[]"],
    ["AsiaWindowsCashtillTool_Ver", "1.0.14", "pp", "[]"],
    ["AsiaWindowsCashtillTool_Ver", "1.0.12", "pr", "[]"],
    ["Test_Store_List", "1277", "pp", "[]"],
]

# ==================== CSV 文件定义 ====================

CSV_FILES = {
    "Module.csv": {
        "headers": ["Name", "Location", "Modules", "Env", "ID"],
        "data_key": "modules",
        "sample": SAMPLE_MODULES,
    },
    "Rule.csv": {
        "headers": ["Name", "Type", "Location", "ENV_NAME", "Rule", "Result"],
        "data_key": "rules",
        "sample": SAMPLE_RULES,
    },
    "Store.csv": {
        "headers": ["Store_Number", "Network_Segment", "Webpos_Env", "EFT", "Location", "RF_Server", "Cashtill_Seg_GW"],
        "data_key": "stores",
        "sample": SAMPLE_STORES,
    },
    "TillList.csv": {
        "headers": ["HostName", "MacAddress"],
        "data_key": "tills",
        "sample": SAMPLE_TILLS,
    },
    "Var.csv": {
        "headers": ["Var_Name", "Value", "Env", "Matcher"],
        "data_key": "vars",
        "sample": SAMPLE_VARS,
    },
}


def load_data_from_yaml(yaml_path):
    """从 YAML 文件加载数据"""
    with open(yaml_path, "r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    csv_data = data.get("csv_data", {}) if data else {}
    return csv_data


def write_csv(filepath, headers, rows):
    """写入 CSV 文件，确保完全匹配目标格式"""
    os.makedirs(os.path.dirname(filepath) if os.path.dirname(filepath) else ".", exist_ok=True)
    with open(filepath, "w", newline="", encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow(headers)
        for row in rows:
            writer.writerow(row)
    print(f"  [OK] {filepath} ({len(rows)} rows)")


def load_module_json_field(data):
    """确保 Module 的 modules 字段是双引号 JSON 字符串"""
    result = []
    for row in data:
        row = list(row)
        if len(row) >= 3:
            val = row[2]
            if isinstance(val, (list, dict)):
                json_str = json.dumps(val, ensure_ascii=False)
                row[2] = json_str
            elif isinstance(val, str) and not val.startswith("["):
                # 尝试解析再序列化，确保格式一致
                try:
                    parsed = json.loads(val)
                    row[2] = json.dumps(parsed, ensure_ascii=False)
                except json.JSONDecodeError:
                    pass
        result.append(row)
    return result


def generate_all(output_dir, csv_data=None):
    """生成所有 CSV 文件"""
    if csv_data is None:
        csv_data = {}

    print(f"\n生成 CSV 文件到: {output_dir}\n")

    for filename, config in CSV_FILES.items():
        rows = csv_data.get(config["data_key"], config["sample"])
        filepath = os.path.join(output_dir, filename)

        if filename == "Module.csv":
            rows = load_module_json_field(rows)

        write_csv(filepath, config["headers"], rows)

    print(f"\n全部生成完毕！共 {len(CSV_FILES)} 个文件\n")


def main():
    parser = argparse.ArgumentParser(description="CSV 文件生成器 - 生成收银台配置 CSV 文件")
    parser.add_argument("-i", "--input", help="YAML 数据文件路径（可选，默认使用内置样例数据）")
    parser.add_argument("-o", "--output", default="generated_csv", help="输出目录（默认: generated_csv）")
    args = parser.parse_args()

    csv_data = None
    if args.input:
        if not os.path.exists(args.input):
            print(f"[ERROR] 文件不存在: {args.input}")
            sys.exit(1)
        print(f"从 {args.input} 加载数据...")
        csv_data = load_data_from_yaml(args.input)

    generate_all(args.output, csv_data)


if __name__ == "__main__":
    main()
