<template>
  <q-page class="flex flex-center">
    <q-card class="singin-container">
      <q-card-section class="bg-teal text-white">
        <div class="text-h6 text-center">登入</div>
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
              v-model="loginName"
              label="登入帳號"
              bottom-slots
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
              :type="isPwd ? 'password' : 'text'"
              v-model="password"
              label="密碼"
              bottom-slots
              :rules="[(val) => 6 <= val.length || '密碼長度最少6個字元']"
              @update:model-value="formatVaild"
              ><template v-slot:append>
                <q-icon
                  :name="isPwd ? 'visibility_off' : 'visibility'"
                  class="cursor-pointer"
                  @click="isPwd = !isPwd"
                /> </template
            ></q-input>
          </div>
          <div class="row q-px-md items-center">
            <q-icon
              name="warning"
              class="text-red"
              v-if="loginFailed"
              style="font-size: 15px"
            />
            <div class="text-red-14">&ensp;{{ loginFailedMsg }}&ensp;</div>
          </div>
        </div>
        <div class="row justify-end q-px-sm q-py-none">
          <q-btn
            color="deep-orange"
            glossy
            label="登入"
            :disable="!signinEnable"
            @click="sendSigninData"
          />
        </div>
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<style lang="scss" scoped>
@import '../css/width.scss';

.singin-container {
  @include xs-width {
    width: 95%;
  }

  @media (min-width: $breakpoint-xs) {
    width: 570px;
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
import axios from 'axios';
import { api } from 'boot/axios';
import { Push } from 'src/lib/router/pushPage';

export default defineComponent({
  name: 'RegisterPage',
  setup() {
    const push = new Push();

    const loginName = ref('');
    const password = ref('');
    const isPwd = ref(true);

    const signinEnable = ref(false);
    const loginFailed = ref(false);
    const loginFailedMsg = ref('');

    // 檢查帳號密碼格式
    function formatVaild() {
      loginFailed.value = false;
      loginFailedMsg.value = '';
      const loginNameLen = loginName.value.length;
      // 檢查 login name  長度
      if (4 > loginNameLen || loginNameLen > 30) {
        signinEnable.value = false;
        return false;
      }
      // 檢查 password 長度
      if (password.value.length < 6) {
        signinEnable.value = false;
        return false;
      }
      // 格式皆符合
      signinEnable.value = true;
      return true;
    }

    // 傳送登入資料
    async function sendSigninData() {
      signinEnable.value = false;
      if (formatVaild()) {
        let path = '/signin';
        await api
          .post(path, {
            login_name: loginName.value,
            password: password.value,
          })
          .then(() => {
            // 成功登入 返回上一頁
            push.previousPage();
          })
          .catch((error) => {
            if (axios.isAxiosError(error)) {
              if (error.response) {
                signinEnable.value = false;
                switch (error.response.status) {
                  case 401:
                    loginFailed.value = true;
                    loginFailedMsg.value = '帳號或密碼錯誤';
                    break;
                  default:
                    console.log('other error', error.response);
                }
              }
            }
          });
      } else {
        signinEnable.value = false;
      }
    }

    return {
      loginName,
      password,
      isPwd,

      signinEnable,
      loginFailed,
      loginFailedMsg,

      formatVaild,
      sendSigninData,
    };
  },
  methods: {},
});
</script>
