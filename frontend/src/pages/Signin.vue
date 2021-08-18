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
              lazy-rules
              :rules="[
                val => 4 <= val.length || '長度需要在4~30個字元之間',
                val => val.length <= 30 || '長度需要在4~30個字元之間'
              ]"
              @input="formatVaild"
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
              lazy-rules
              :rules="[val => 6 <= val.length || '密碼長度最少6個字元']"
              @input="formatVaild"
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
              style="font-size: 15px;"
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
@import "../css/width.scss";

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

<script>
export default {
  name: "RegisterPage",
  data() {
    return {
      loginName: "",
      password: "",
      isPwd: true,

      signinEnable: false,
      loginFailed: false,
      loginFailedMsg: ""
    };
  },
  computed: {},
  methods: {
    formatVaild() {
      this.loginFailed = false;
      this.loginFailedMsg = "";
      if (4 > this.loginName || this.loginName.length > 30) {
        this.signinEnable = false;
        return false;
      }

      if (this.password.length < 6) {
        this.signinEnable = false;
        return false;
      }

      this.signinEnable = true;
      return true;
    },
    redirectSigninPage() {
      this.$router.push("/signin");
    },
    async sendSigninData() {
      this.signinEnable = false;
      if (this.formatVaild()) {
        var path = "/api/signin";
        this.$axios
          .post(path, {
            login_name: this.loginName,
            password: this.password
          })
          .then(res => {
            var data = res.data;
            this.id = data["user_id"];
            this.errList = data["show_name"];
            this.$router.go(-1);
            console.log("成功登入");
          })
          .catch(error => {
            if (error.response) {
              this.signinEnable = false;
              // 當狀態碼不在 validateStatus 設定的範圍時進入
              // 有 data / status / headers 參數可用
              switch (error.response.status) {
                case 401:
                  this.loginFailed = true;
                  this.loginFailedMsg = "帳號或密碼錯誤";
                  break;
                default:
                  console.log("other error", error.response);
              }
            } else if (error.request) {
              // 發送請求，但沒有接到回應
              // error.request is an instance of XMLHttpRequest in the browser
              // and an instance of http.ClientRequest in node.js
              console.log(error.request);
            } else {
              // 在設定 request 時出錯會進入此
              console.log("Error", error.message);
            }
          });
      } else {
        this.signinEnable = false;
      }
    }
  }
};
</script>
