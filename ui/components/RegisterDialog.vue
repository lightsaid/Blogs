<template>
    <button type="button" class="btn btn-outline-light ms-3 fw-semibold" data-bs-toggle="modal"
        data-bs-target="#registerModal"
        style="--bs-btn-padding-y: .3rem; --bs-btn-padding-x: .5rem; --bs-btn-font-size: .85rem;">
        Register
    </button>

    <!-- Modal -->
    <div class="modal fade" style="--bs-backdrop-zindex=7" id="registerModal" tabindex="-1"
        aria-labelledby="registerModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="registerModalLabel">Register</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <form class="needs-validation" novalidate>
                        <div class="mb-3">
                            <label for="username" class="col-form-label">UserName:</label>
                            <input type="text" class="form-control is-invalid" id="username" name="username" v-model="data.username" required />
                            <div class="invalid-feedback">请输入用户名</div>
                        </div>
                        <div class="mb-3">
                            <label for="email" class="col-form-label">Email:</label>
                            <input type="email" class="form-control" id="email" name="email" v-model="data.email" required />
                            <div class="invalid-feedback">请输入邮箱地址</div>
                        </div>
                        <div class="mb-3">
                            <label for="password" class="col-form-label">Password:</label>
                            <input type="password" class="form-control" id="password" name="password" v-model="data.password" required />
                            <div class="invalid-feedback">请输入密码</div>
                        </div>
                        <div class="mb-3">
                            <label for="repassword" class="col-form-label">Repassword:</label>
                            <input type="password" class="form-control" id="repassword" name="repassword" v-model="data.repassword" is-valid
                                required />
                            <div class="invalid-feedback">请输入确认密码</div>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-outline-secondary" data-bs-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-primary" @click="registerHandler">Register</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const data = reactive({
    username: "",
    email: "",
    password: "",
    repassword: "",
})

const errors = reactive({
    username: {error:"请输入用户名", validator: [(value:string)=>{}]},
    email: {error:"请输入邮箱地址", validator: [(value:string)=>{}]},
    password: {error:"请输入密码", validator: [(value:string)=>{}]},
    repassword: {error:"请再次输入密码", validator: [(value:string)=>{}]}
})

function validate() {
    const form = document.querySelector('.needs-validation')
    form?.classList.add('was-validated')

    


    return true
}

function registerHandler() {
    // 验证通过
    if (validate()) {
        console.log(data)

        return
    }
}
</script>

<style lang="scss">
.modal-backdrop {
    --bs-backdrop-zindex: 7 !important;
}
</style>