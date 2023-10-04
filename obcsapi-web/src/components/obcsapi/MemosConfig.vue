<script setup lang="ts">
import { ref, onMounted } from "vue";
import VueForm from "@lljj/vue3-form-naive"
import { ObcsapiConfigGet, ObcsapiConfigPost, ObcsapiServerInfo, ObcsapiSetOAuth2 } from "@/api/obcsapi"
import { NScrollbar } from "naive-ui"
import { ObcsapiTestMail, ObcsapiUpdateBdGet, ObcsapiUpdateConfig } from "@/api/obcsapi";


const formData = ref({});
const schema = ref({
    type: "object",
    properties: {
        basic: {
            title: "基础设置",
            type: "object",
            properties: {
                disable_login: {
                    title: "禁用账户密码登录",
                    type: "boolean",
                    description: "禁用登录，已经下发的 token 可以继续使用。",
                    'ui:options': {
                        placeholder: false,
                    }
                }
            }
        },
        ob_daily_config: {
            title: "Obsidian Daily 设置",
            type: "object",
            description: "日记及附件存放位置。日期格式化 2006-01-02 15:04:05",
            properties: {
                ob_daily_dir: {
                    type: "string",
                    title: "日记文件夹",
                    'ui:options': {
                        placeholder: "日记/",
                    }
                },
                ob_daily: {
                    type: "string",
                    title: "日记文件格式",
                    description: "格式化时间",
                    'ui:options': {
                        placeholder: "200601/2006-01-02",
                    }
                },
                ob_daily_attachment_dir_under_daily: {
                    type: "boolean",
                    title: "日记附件文件夹前缀",
                    description: "日记附件文件夹: true 根目录下 false: 日记文件夹下",
                },
                ob_daily_attachment_dir: {
                    type: "string",
                    title: "日记附件文件夹",
                    description: "日记附件文件夹 格式化时间",
                    'ui:options': {
                        placeholder: "附件/200601/",
                    }
                },
                ob_other_data_dir: {
                    type: "string",
                    title: "其他文件夹",
                    'ui:options': {
                        placeholder: "其他文件/",
                    }
                },

            }
        },
        wechat_mp: {
            title: "微信公众号",
            type: "object",
            properties: {
                return_str: {
                    type: "string",
                    title: "返回字符串",
                    'ui:options': {
                        placeholder: "📩 已保存，\u003ca href='https://kkbt.gitee.io/web/'\u003e点击查看今日笔记\u003c/a\u003e",
                    }
                }
            }
        },
        webdav: {
            title: "WebDAV",
            type: "object",
            description: "LocalStorage (RemotelySave WebDav) 用户自定义账户密码",
            properties: {
                server: {
                    type: "boolean",
                    title: "服务开关",
                },
                username: {
                    type: "string",
                    title: "WebDAV 自定义用户名",
                    'ui:options': {
                        placeholder: "kkbt",
                    }
                },
                password: {
                    type: "string",
                    title: "WebDAV 自定义密码",
                    'ui:options': {
                        placeholder: "webdavpassword",
                    }
                },
                ob_local_dir: {
                    type: "string",
                    title: "库文件夹位置",
                    description: "数据源选择本地时，存放的本地文件夹位置。需要和removely save文件夹一样，正常为 Ob 库的名。",
                    'ui:options': {
                        placeholder: "日记/",
                    }
                }
            }
        },
        mail: {
            title: "邮件服务",
            type: "object",
            description: "用于提醒服务",
            properties: {
                smtp_host: {
                    type: "string",
                    title: "smtp_host",
                    'ui:options': {
                        placeholder: "smtpdm.aliyun.com",
                    }
                },
                smtp_port: {
                    type: "number",
                    title: "smtp_port",
                    'ui:options': {
                        placeholder: "80",
                    }
                },
                user_name: {
                    type: "string",
                    title: "账户",
                    'ui:options': {
                        placeholder: "no-reply@example.com",
                    }
                },
                password: {
                    type: "string",
                    title: "密码",
                    'ui:options': {
                        placeholder: "xxxxxxxx",
                    }
                },
                sender_email: {
                    type: "string",
                    title: "发送者邮箱",
                    'ui:options': {
                        placeholder: "no-reply@example.com",
                    }
                },
                sender_name: {
                    type: "string",
                    title: "发送者名",
                    'ui:options': {
                        placeholder: "ObCSAPI",
                    }
                },
                receiver_email: {
                    type: "string",
                    title: "接收者邮箱",
                    'ui:options': {
                        placeholder: "xxx@gmail.com",
                    }
                },
            }
        },
        image_hosting: {
            title: "图床",
            type: "object",
            description: "ImageHosting 图床文件 有四部分构成 url 文件夹及前缀，原名字，随机字符。图床文件夹及文件前缀 eg 2006-01-02 15:04:05 如 按月存放是 01/ ; 按 年存放 2006/ ; 文件前缀 200601 ; 文件夹和文件前缀 200601/200601_",
            properties: {
                storage_mode: {
                    type: "string",
                    title: "存储位置选择",
                    description: "local 服务器存储；obsidian Ob库存储；s3 对象存储",
                    enum: ["local", "obsidian", "s3"],
                    enumNames: ["local", "obsidian", "s3"]
                },
                base_url: {
                    type: "string",
                    title: "BaseUrl ",
                    'ui:options': {
                        placeholder: "http://localhost:8900/images/",
                    }
                },
                prefix: {
                    type: "string",
                    title: "时间格式化 Prefix",
                    description: "时间格式化",
                    'ui:options': {
                        placeholder: "200601/kkbt_",
                    }
                },
                use_raw_name: {
                    type: "boolean",
                    title: "是否保留文件原名",
                },
                random_char_length: {
                    type: "number",
                    title: "随机字符串长度",
                    'ui:options': {
                        placeholder: "5",
                    }
                },
                use_bd_ocr: {
                    type: "boolean",
                    title: "是否使用 百度 OCR",
                },
                bd_ocr_access_token: {
                    type: "string",
                    title: "百度 OCR Access Token",
                    'ui:options': {
                        placeholder: "xxxxx.xxxxx.xxxxx.xxxxx.xxxxx-xxxxx",
                    }
                },
            }
        },
        s3_compatible: {
            title: "S3 兼容存储",
            type: "object",
            description: "上传图床时使用的 S3 配置",
            properties: {
                end_point: {
                    type: "string",
                    description: "End Point",
                    'ui:options': {
                        placeholder: "https://s3-cn-south-1.qiniucs.com",
                    }
                },
                region: {
                    type: "string",
                    description: "Region",
                    'ui:options': {
                        placeholder: "s3-cn-south-1",
                    }
                },
                bucket: {
                    type: "string",
                    description: "Bucket",
                    'ui:options': {
                        placeholder: "bucketname",
                    }
                },
                access_key: {
                    type: "string",
                    description: "Access Key",
                    'ui:options': {
                        placeholder: "xxxxx",
                    }
                },
                secret_key: {
                    "type": "string",
                    "description": "Secret Key",
                    'ui:options': {
                        placeholder: "xxxxx",
                    }
                },
                base_url: {
                    "type": "string",
                    "description": "自定义域名，此项可以空，但图片可能会访问错误",
                    'ui:options': {
                        placeholder: "https://youdomain.com",
                    }
                }
            }
        },
        bd_ocr: {
            title: "百度OCR",
            type: "object",
            description: "百度 OCR 配置",
            properties: {
                api_key: {
                    type: "string",
                    title: "Api Key",
                    'ui:options': {
                        placeholder: "xxxxx",
                    }
                },
                api_secret: {
                    type: "string",
                    title: "Api Secret",
                    'ui:options': {
                        placeholder: "xxxxx",
                    }
                }
            }
        },
        reminder: {
            title: "提醒服务",
            type: "object",
            properties: {
                daily_email_remder_time: {
                    type: "string",
                    title: "每日提醒时间",
                    'ui:options': {
                        placeholder: "0800",
                    }
                },
                reminder_dicionary: {
                    type: "string",
                    title: "微信识别时间所用字典",
                    description: "可选 full  200k 100k 20k  10k",
                    'ui:options': {
                        placeholder: "dictionary-200k.txt",
                    }
                },
            }
        },
        mention: {
            type: "object",
            properties: {
                tags: {
                    type: "array",
                    items: {
                        type: "string",
                    }
                }
            },
            "ui:hidden": true,
        },
        oauth2userinfo: {
            title: "第三方登录设置",
            type: "object",
            properties: {
                gitee_user_info: {
                    title: "Gitee 账户",
                    description: "该项目关联后自动生成",
                    type: "object",
                    properties: {
                        is_active: {
                            type: "boolean",
                            title: "激活"
                        },
                        id: {
                            type: "number",
                        },
                        login: {
                            type: "string",
                        },
                        name: {
                            type: "string",
                        }
                    }
                }
            }
        },
    }
});
const info = ref();
const showInfo = ref(false);
onMounted(() => {
    getConfiguration()
    ObcsapiServerInfo().then(res => {
        info.value = res;
    })
})

function getConfiguration() {
    ObcsapiConfigGet().then(data => {
        formData.value = data
    }).catch((e) => {
        window.$message.error(e.message)
    });
}

function handlerSubmit() {
    ObcsapiConfigPost(formData.value).then(res => {
        window.$message.success(res)
        getConfiguration()
    }).catch(e => {
        window.$message.error(e)
    })
}

function sendMail() {
    ObcsapiTestMail().then(res => {
        window.$message.info(res)
    }).catch(e => {
        window.$message.info(e)
    })
}

function obcsapiUpdateConfig() {
    // update  config.yaml
    ObcsapiUpdateConfig().then(res => {
        window.$message.info(res)
    }).catch(e => {
        window.$message.info(e)
    })
}

function upDateBdOcrAccessToken() {
    ObcsapiUpdateBdGet().then(res => {
        window.$message.info(JSON.stringify(res))
        getConfiguration()
    }).catch(e => {
        window.$message.info(e)
    })
}

function setOAuth2() {
    ObcsapiSetOAuth2().then(url => {
        window.location.href = url;
    }).catch(e => {
        console.log(e);
    });
}

</script>
<template>
    <h1 @click="showInfo = !showInfo"><a>Server Setting</a></h1>
    <div v-if="showInfo && info">
        <a href="https://gitee.com/kkbt/obcsapi-go">Obsidian 云存储后端 API Go 版本项目地址 </a>
        <a href="https://kkbt.gitee.io/obcsapi-go/#/"> 📄文档</a><br>
        ServerTime: {{ info.server_time }}<br>
        ServerVersion: {{ info.server_version }}<br>
        ServerConfigVersion: {{ info.config_version }} <br>
        <a>注意：服务器不支持设置为空，所以如果想置空某一项，只能在服务器处修改 config.run.json 文件 , 然后重启程序</a>
    </div>
    <n-scrollbar style="max-height: 75vh">
        <vue-form v-model="formData" :schema="schema" @submit="handlerSubmit" />
        <n-button @click="sendMail" quaternary>测试邮件</n-button>
        <n-button @click="upDateBdOcrAccessToken" quaternary>更新BD OCR</n-button>
        <n-button @click="obcsapiUpdateConfig" quaternary>更新config.yaml</n-button>
        <n-button @click="setOAuth2" quaternary>关联 Gitee</n-button>
    </n-scrollbar>
    <n-button @click="handlerSubmit" quaternary>保存</n-button>
</template>