<script lang="ts">
import { ObcsapiTalk } from "@/api/obcsapi";
import { NButton } from "naive-ui";
import { defineComponent, h, ref, type Ref } from "vue";
const componentsList: Ref<Array<any>> = ref([]);

function handleClick(text: string) {
    console.log(text);
    SendCommand(text);
}


let obj = {
    "code": 200,
    "data": {
        "type": "message",
        "parts": [
            { "type": "md", "content": "xxxx" },
            { "type": "button", "content": "button1", "command": "#obcsapi-command-json This is button content1" },
            { "type": "button", "content": "button2", "command": "#obcsapi-command-json This is button content2" }
        ]
    },
    "msg": "This is a message"
}

RenderParts(obj);


function RenderParts(obj: any) {
    let parts = obj.data.parts;
    for (let i = 0; i < parts.length; i++) {
        let part = parts[i];
        switch (part.type) {
            case "md":
                componentsList.value.push(h("div", part.content))
                break;
            case "button":
                componentsList.value.push(
                    h(NButton, { onClick: () => handleClick(part.command) }, () => part.content)
                )
        }
    }
}

function SendCommand(command:string) {
    ObcsapiTalk(command).then(text => {
        console.log(text);
    }).catch(e => {
        console.log(e);
        window.$message.warning("Err Save: " + e);
    })
}


export default defineComponent({
    render() {
        return componentsList.value;
    },
});
</script>
