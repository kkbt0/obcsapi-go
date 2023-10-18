<script lang="ts">
import { ObcsapiPostMemos, ObcsapiTalk } from "@/api/obcsapi";
import { NButton } from "naive-ui";
import { defineComponent, h, ref, render, type Ref } from "vue";
import TBB from "@/components/talkpage/TBB.vue";
import Basic from "@/components/talkpage/Basic.vue";
import { Adjutant } from "@/stores/adjutant";
// example
// let obj = {
//     "code": 200,
//     "data": {
//         "type": "message",
//         "command": "#xxx xxx",
//         "parts": [
//             { "type": "md", "text": "xxxx" },
//             { "type": "button", "text": "button1", "command": "#obcsapi-command-json This is button content1" },
//             { "type": "button", "text": "button2", "command": "#obcsapi-command-json This is button content2" }
//         ]
//     },
//     "msg": "This is a message"
// }

export default defineComponent({
    props: {
        inText: {
            type: String,
            required: true
        },
    },
    setup(props) {
        const componentsList: Ref<Array<any>> = ref([]);

        TryRender(props.inText);

        function TryRender(text: string) {
            try {
                componentsList.value = [];
                let obj = JSON.parse(text);
                RenderParts(obj);
            } catch (error) {
                componentsList.value = [];
                componentsList.value.push(h(Basic, { mdText: text }));
            }

        }
        function handleClick(obj:any,text: string) {
            console.log(text);
            SendCommand(obj,text);
        }

        function RenderParts(obj: any) {
            if (obj.data.type == "message-tbb") {
                componentsList.value.push(h(TBB, {
                    inobj: obj,
                    "onWeb-command": SendCommand
                }))
                return
            }
            let parts = obj.data.parts;
            for (let i = 0; i < parts.length; i++) {
                let part = parts[i];
                switch (part.type) {
                    case "md":
                        componentsList.value.push(h(Basic, { mdText: part.text }));
                        break;
                    case "button":
                        componentsList.value.push(
                            h(NButton, { onClick: () => handleClick(obj,part.command) }, () => part.text)
                        )
                }
            }
        }

        function SendCommand(obj:any,command: string) {
            if (command.startsWith("@")) {
                switch (command) {
                    case "@web-save-memos":
                        ObcsapiPostMemos("", 9999, "", obj.data.parts[0].text).then(data => {
                            Adjutant().success("Successfully Sent")
                        }).catch(e => {
                            console.log(e);
                        });
                        break;
                    case "@web-resend":
                        ObcsapiTalk(obj.data.command).then(text => {
                            TryRender(text);
                        }).catch(e => {
                            console.log(e);
                        })
                        break;
                }
            } else {
                ObcsapiTalk(command).then(text => {
                    TryRender(text);
                }).catch(e => {
                    console.log(e);
                })
            }

        }

        return {
            componentsList
        }
    },
    render() {
        return this.componentsList;
    },
});
</script>
