<template>
  <q-page class="flex flex-center">
    <q-dialog v-model="registerError">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">註冊失敗</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <p v-for="msg in registerErrMsgArray" :key="msg">
            {{ msg }}
          </p>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="確認" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="registerSuccess">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">註冊成功</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          {{ registerSuccessMsg }}
        </q-card-section>

        <q-card-actions align="right">
          <q-btn
            flat
            label="返回首頁"
            color="primary"
            v-close-popup
            @click="push.homePage()"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-card class="register-container">
      <q-card-section class="bg-teal text-white">
        <div class="text-h6 text-center">註冊帳號</div>
        <div class="text-subtitle2"></div>
      </q-card-section>

      <q-card-actions vertical align="center">
        <div class="column input-container">
          <div class="column q-pa-sm">
            <q-input
              class="input"
              outlined
              dense
              bg-color="grey-4"
              v-model="showName"
              label="顯示名稱"
              bottom-slots
              hint="公開用的名稱，註冊後仍然可以更改"
              :rules="[
                (val) => 1 <= val.length || '長度需要在1~15個字元之間',
                (val) => val.length <= 15 || '長度需要在1~15個字元之間',
              ]"
              @update:model-value="formatVaild"
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              outlined
              dense
              bg-color="grey-4"
              v-model="loginName"
              label="登入帳號"
              bottom-slots
              hint="登入用的帳號"
              :rules="[
                (val) => 4 <= val.length || '長度需要在4~30個字元之間',
                (val) => val.length <= 30 || '長度需要在4~30個字元之間',
              ]"
              @update:model-value="formatVaild"
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              outlined
              dense
              bg-color="grey-4"
              type="email"
              v-model="email"
              label="電子郵件"
              bottom-slots
              hint="用來修改密碼或是找回密碼"
              :rules="[
                (val) => val.length <= 256 || '不支援長度超過256字元的電子郵件',
                (val) => val.search(this.emailRule) != -1 || '電子郵件格式錯誤',
              ]"
              @update:model-value="formatVaild"
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              ref="passwordInputRef"
              outlined
              dense
              bg-color="grey-4"
              v-model="password"
              type="password"
              label="密碼"
              bottom-slots
              hint="長度最小6個字元"
              :rules="[
                (val) => 6 <= val.length || '密碼長度最小6個字元',
                (val) => val == this.rePassword || '密碼不一樣',
              ]"
              @update:model-value="
                formatVaild();
                passwordVaild();
              "
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              ref="rePasswordInputRef"
              outlined
              dense
              bg-color="grey-4"
              v-model="rePassword"
              type="password"
              label="確認密碼"
              :rules="[
                (val) => 6 <= val.length || '密碼長度最小6個字元',
                (val) => val == this.password || '密碼不一樣',
              ]"
              @update:model-value="
                formatVaild();
                passwordVaild();
              "
            />
          </div>
        </div>
        <div class="row justify-between q-px-md q-py-sm">
          <q-btn
            color="white"
            text-color="black"
            label="登入頁面"
            @click="push.signinPage()"
          />
          <q-btn
            color="deep-orange"
            glossy
            label="註冊"
            :disable="!registerEnable"
            @click="sendRegisterData"
          />
        </div>
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<style lang="scss" scoped>
@import '../css/width.scss';

.register-container {
  @include xs-width {
    width: 95%;
  }

  @media (min-width: $breakpoint-xs) {
    width: 538px;
  }
}

.input-container {
  @include xs-width {
    width: 100%;
  }

  @media (min-width: $breakpoint-xs) {
    width: 100%;
  }
}
.input {
  @include xs-width {
    width: 100%;
  }

  @media (min-width: $breakpoint-xs) {
    width: 100%;
  }
}
</style>

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { QInput } from 'quasar';
import { api } from 'boot/axios';
import { Push } from 'src/lib/router/pushPage';

export default defineComponent({
  name: 'RegisterPage',
  setup() {
    const push = new Push();

    const requireEmailAct = ref(false);
    const showName = ref('');
    const loginName = ref('');
    const email = ref('');
    const emailRule =
      /^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z]+$/;
    const password = ref('');
    const passwordInputRef = ref(null);
    const rePassword = ref('');
    const rePasswordInputRef = ref(null);
    const registerEnable = ref(false);
    const registerError = ref(false);
    const registerErrMsgArray = ref([] as string[]);
    const registerSuccess = ref(false);
    const registerSuccessMsg = ref('');
    const id = ref(-1);
    const errList = ref([]);

    interface serverInfoDataJson {
      require_email_activate: boolean;
    }

    async function getServerInfo() {
      let path = '/server-info';
      await api.get(path).then((res) => {
        let data = res.data as serverInfoDataJson;
        requireEmailAct.value = data['require_email_activate'];
      });
    }
    void getServerInfo();

    function passwordEqual() {
      if (password.value == rePassword.value) {
        return true;
      } else {
        return false;
      }
    }

    function passwordVaild() {
      if (passwordEqual() && password.value.length >= 6) {
        (passwordInputRef.value as unknown as QInput).resetValidation();
        (rePasswordInputRef.value as unknown as QInput).resetValidation();
      }
    }

    // formatVaild 檢測所有輸入格式是否正確
    function formatVaild() {
      if (1 > showName.value.length || showName.value.length > 15) {
        registerEnable.value = false;
        return false;
      }
      if (4 > loginName.value.length || loginName.value.length > 30) {
        registerEnable.value = false;
        return false;
      }

      if (email.value.length > 256 || email.value.search(emailRule) == -1) {
        registerEnable.value = false;
        return false;
      }
      if (!passwordEqual() || password.value.length < 6) {
        registerEnable.value = false;
        return false;
      }
      registerEnable.value = true;
      return true;
    }

    async function sendRegisterData() {
      if (formatVaild()) {
        var path = '/api/register';
        await api
          .post(path, {
            show_name: showName.value,
            login_name: loginName.value,
            email: email.value,
            password: password.value,
          })
          .then((res) => {
            interface resDataJson {
              id: number;
              err_list: [];
            }
            var data = res.data as resDataJson;
            id.value = data['id'];
            errList.value = data['err_list'];
            if (errList.value) {
              registerErrMsgArray.value = [] as string[];
              errList.value.forEach((errCode) => {
                switch (errCode) {
                  case 1:
                    registerErrMsgArray.value.push('登入帳號已被使用過');
                    break;
                  case 2:
                    registerErrMsgArray.value.push('帳號長度錯誤');
                    break;
                  case 3:
                    registerErrMsgArray.value.push('顯示名稱長度錯誤');

                    break;
                  case 4:
                    registerErrMsgArray.value.push('密碼長度錯誤');
                    break;
                  case 5:
                    registerErrMsgArray.value.push(
                      '不支援大於256字元的電子郵件'
                    );

                    break;
                  case 6:
                    registerErrMsgArray.value.push('電子郵件格式錯誤');
                    break;
                  case 7:
                    registerErrMsgArray.value.push('電子郵件已被使用過');
                    break;
                  case 8:
                    registerErrMsgArray.value.push('系統錯誤');
                    break;
                  default:
                    break;
                }
              });
              registerError.value = true;
            } else {
              if (requireEmailAct.value) {
                registerSuccessMsg.value =
                  '認證信件已經寄到 ' +
                  email.value +
                  ' 請透過信中的連結完成註冊';
              }
              registerSuccess.value = true;
            }
          })
          .catch((error) => {
            if (error) {
            }
          });
      } else {
      }
    }

    return {
      requireEmailAct,
      showName,
      loginName,
      email,
      emailRule,
      password,
      rePassword,
      passwordInputRef,
      rePasswordInputRef,
      registerEnable,
      registerError,
      registerErrMsgArray,
      registerSuccess,
      registerSuccessMsg,
      id,
      errList,
      push,
      passwordVaild,
      formatVaild,
      sendRegisterData,
    };
  },
  methods: {},
});
</script>
