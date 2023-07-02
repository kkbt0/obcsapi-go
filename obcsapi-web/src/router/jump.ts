import { ObcsapiLoginByCode, ObcsapiTestJwt } from "@/api/obcsapi";
import { useRoute, useRouter } from 'vue-router'

export const LoginByCode = () => {
    let code = useRoute().query['code']?.toString();
    let router = useRouter();
    console.log("Code",code);
    if (code == null || code == undefined || code == "") {
        return
    }
    ObcsapiLoginByCode(code).then(obj => {
        if (obj.code == 200 || obj.code == 201) {
            localStorage.setItem("token", obj.token)
            window.$message.success("登录成功")
            router.push("/").then(() => {
                location. reload();
            })
        } else {
            window.$message.warning("Failed to login");
        }
    });
}

export const SayHello = () => {
    if (window.location.href.includes("code=")) {
        return 
    }
    console.log(`Initiating Hello`)
    const router = useRouter();
    ObcsapiTestJwt().then(res => {
        if (res.code != 200) {
            router.push("/login").then(() => {
                location.reload();
            })
            window.$message.error(`${res.code}`);
        }
    }).catch(err => {
        window.$message.error(err);
    })
}
