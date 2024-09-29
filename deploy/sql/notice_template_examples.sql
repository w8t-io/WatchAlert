use watchalert;
INSERT ignore  INTO `notice_template_examples` (`id`, `name`, `description`, `template`, `enable_fei_shu_json_card`,
                                        `template_firing`, `template_recover`, `notice_type`)
VALUES ('nt-cqh3uppd6gvj2ctaqd60', '通用模版', '', '{{- define "Title" -}}
{{- if not .IsRecovered -}}
    【报警中】- 即时设计业务系统 🔥
{{- else -}}
    【已恢复】- 即时设计业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
**🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.instance}
    **🕘 开始时间:** ${first_trigger_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations}
    {{- else -}}
    **🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.instance}
    **🕘 开始时间:** ${first_trigger_time_format}
    **🕘 恢复时间:** ${recover_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations}
    {{- end -}}
    {{ end }}

    {{- define "Footer" -}}
    🧑‍💻 即时设计 - 运维团队
{{- end }}', false, '', '', 'FeiShu');
INSERT  ignore INTO `notice_template_examples` (`id`, `name`, `description`, `template`, `enable_fei_shu_json_card`,
                                        `template_firing`, `template_recover`, `notice_type`)
VALUES ('nt-cqh4361d6gvj80netqk0', '飞书高级消息卡片模版', '', '', true, '{
  "elements": [
    {
      "tag": "column_set",
        "flexMode": "none",
        "background_style": "default",
        "text": {
        "content": "",
        "tag": ""
      },
        "actions": null,
        "columns": [],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**🫧 报警指纹：**\n${fingerprint}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🤖 报警类型：**\n${rule_name}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**📌 报警等级：**\n${severity}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🕘 开始时间：**\n${first_trigger_time_format}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**👤 值班人员：**\n${duty_user}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🖥 报警主机：**\n${metric.instance}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**📝 报警事件：**\n${annotations}",
        "tag": "lark_md"
              }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "hr",
            "flexMode": "",
                "background_style": "",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": null,
        "elements": null
    },
        { "tag" : "note",
            "flexMode": "",
                "background_style": "",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": null,
        "elements": [
        {
          "tag": "plain_text",
        "content": "🧑‍💻 即时设计 - 运维团队"
        }
      ]
    }
  ],
        "header": {
    "template": "red",
        "title": {
      "content": "【报警中】- 即时设计业务系统 🔥",
        "tag": "plain_text"
    }
  },
        "tag": ""
}', '{
  "elements": [
    {
      "tag": "column_set",
        "flexMode": "none",
        "background_style": "default",
        "text": {
        "content": "",
        "tag": ""
      },
        "actions": null,
        "columns": [],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**🫧 报警指纹：**\n${fingerprint}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🤖 报警类型：**\n${rule_name}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**📌 报警等级：**\n${severity}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🕘 开始时间：**\n${first_trigger_time_format}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**🕘 恢复时间：**\n${recover_time_format}",
        "tag": "lark_md"
              }
            }
          ]
        },
        { "tag" : "column",
            "width": "weighted",
                "weight": 1,
                "vertical_align": "top",
            "elements": [
            {
            "tag": "div",
            "text": {
                "content": "**🖥 报警主机：**\n${metric.instance}",
                "tag": "lark_md"
            }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**👤 值班人员：**\n${duty_user}",
        "tag": "lark_md"
              }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "column_set",
            "flexMode": "none",
                "background_style": "default",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": [
        {
          "tag": "column",
        "width": "weighted",
        "weight": 1,
        "vertical_align": "top",
        "elements": [
            {
              "tag": "div",
        "text": {
                "content": "**📝 报警事件：**\n${annotations}",
        "tag": "lark_md"
              }
            }
          ]
        }
      ],
        "elements": null
    },
        { "tag" : "hr",
            "flexMode": "",
                "background_style": "",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": null,
        "elements": null
    },
        { "tag" : "note",
            "flexMode": "",
                "background_style": "",
                "text": {
            "content": "",
            "tag": ""
            },
        "actions": null,
        "columns": null,
        "elements": [
        {
          "tag": "plain_text",
        "content": "🧑‍💻 即时设计 - 运维团队"
        }
      ]
    }
  ],
        "header": {
    "template": "green",
        "title": {
      "content": "【已恢复】- 即时设计业务系统 ✨",
        "tag": "plain_text"
    }
  },
        "tag": ""
}', 'FeiShu');
INSERT  ignore INTO `notice_template_examples` (`id`, `name`, `description`, `template`, `enable_fei_shu_json_card`,
                                        `template_firing`, `template_recover`, `notice_type`)
VALUES ('nt-cqh4599d6gvj80netql0', 'Email邮件通知模版', '', '{{ define "Event" -}}
{{- if not .IsRecovered -}}
<p>==========<strong>告警通知</strong>==========</p>
<strong>🤖 报警类型:</strong> ${rule_name}<br>
    <strong>🫧 报警指纹:</strong> ${fingerprint}<br>
    <strong>📌 报警等级:</strong> ${severity}<br>
    <strong>🖥 报警主机:</strong> ${metric.node_name}<br>
    <strong>🧚 容器名称:</strong> ${metric.pod}<br>
    <strong>☘️ 业务环境:</strong> ${metric.namespace}<br>
    <strong>🕘 开始时间:</strong> ${first_trigger_time_format}<br>
    <strong>👤 值班人员:</strong> ${duty_user}<br>
    <strong>📝 报警事件:</strong> ${annotations}<br>
    {{- else -}}
    <p>==========<strong>恢复通知</strong>==========</p>
    <strong>🤖 报警类型:</strong> ${rule_name}<br>
    <strong>🫧 报警指纹:</strong> ${fingerprint}<br>
    <strong>📌 报警等级:</strong> ${severity}<br>
    <strong>🖥 报警主机:</strong> ${metric.node_name}<br>
    <strong>🧚 容器名称:</strong> ${metric.pod}<br>
    <strong>☘️ 业务环境:</strong> ${metric.namespace}<br>
    <strong>🕘 开始时间:</strong> ${first_trigger_time_format}<br>
    <strong>🕘 恢复时间:</strong> ${recover_time_format}<br>
    <strong>👤 值班人员:</strong> ${duty_user}<br>
    <strong>📝 报警事件:</strong> ${annotations}<br>
    {{- end -}}
    {{ end }}', false, '', '', 'Email');
INSERT  ignore INTO `notice_template_examples` (`id`, `name`, `description`, `template`, `enable_fei_shu_json_card`,
                                        `template_firing`, `template_recover`, `notice_type`)
VALUES ('nt-cqh45t9d6gvj80netqm0', 'Loki日志告警通知模版', '', '{{- define "Title" -}}
{{- if not .IsRecovered -}}
    【报警中】- 即时设计业务系统 🔥
{{- else -}}
    【已恢复】- 即时设计业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
**🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.node_name}
    **🧚 容器名称:** ${metric.pod}
    **☘️ 业务环境:** ${metric.namespace}
    **🕘 开始时间:** ${first_trigger_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations}
    {{- else -}}
    **🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.node_name}
    **🧚 容器名称:** ${metric.pod}
    **☘️ 业务环境:** ${metric.namespace}
    **🕘 开始时间:** ${first_trigger_time_format}
    **🕘 恢复时间:** ${recover_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations}
    {{- end -}}
    {{ end }}

    {{- define "Footer" -}}
    🧑‍💻 即时设计 - 运维团队
{{- end }}', false, '', '', 'FeiShu');
INSERT  ignore INTO `notice_template_examples` (`id`, `name`, `description`, `template`, `enable_fei_shu_json_card`,
                                        `template_firing`, `template_recover`, `notice_type`)
VALUES ('nt-cqh464hd6gvj80netqng', '阿里云SLS日志告警通知模版', '', '{{- define "Title" -}}
{{- if not .IsRecovered -}}
    【报警中】- 即时设计业务系统 🔥
{{- else -}}
    【已恢复】- 即时设计业务系统 ✨
{{- end -}}
{{- end }}

{{- define "TitleColor" -}}
{{- if not .IsRecovered -}}
red
{{- else -}}
green
{{- end -}}
{{- end }}

{{ define "Event" -}}
{{- if not .IsRecovered -}}
**🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.__tag__:_node_name_}
    **🧚 容器名称:** ${metric._pod_name_}
    **☘️ 业务环境:** ${metric._namespace_}
    **🕘 开始时间:** ${first_trigger_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations.content}
    {{- else -}}
    **🤖 报警类型:** ${rule_name}
    **🫧 报警指纹:** ${fingerprint}
    **📌 报警等级:** ${severity}
    **🖥 报警主机:** ${metric.__tag__:_node_name_}
    **🧚 容器名称:** ${metric._pod_name_}
    **☘️ 业务环境:** ${metric._namespace_}
    **🕘 开始时间:** ${first_trigger_time_format}
    **🕘 恢复时间:** ${recover_time_format}
    **👤 值班人员:** ${duty_user}
    **📝 报警事件:** ${annotations.content}
    {{- end -}}
    {{ end }}

    {{- define "Footer" -}}
    🧑‍💻 即时设计 - 运维团队
{{- end }}', false, '', '', 'FeiShu');