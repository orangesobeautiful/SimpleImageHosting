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
            @click="redirectHomePage"
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
              lazy-rules
              :rules="[
                val => 1 <= val.length || '長度需要在1~15個字元之間',
                val => val.length <= 15 || '長度需要在1~15個字元之間'
              ]"
              @input="formatVaild"
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
              type="email"
              v-model="email"
              label="電子郵件"
              bottom-slots
              hint="用來修改密碼或是找回密碼"
              lazy-rules
              :rules="[
                val => val.length <= 256 || '不支援長度超過256字元的電子郵件',
                val => val.search(this.emailRule) != -1 || '電子郵件格式錯誤'
              ]"
              @input="formatVaild"
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              ref="passwordInput"
              outlined
              dense
              bg-color="grey-4"
              v-model="password"
              type="password"
              label="密碼"
              bottom-slots
              hint="長度最小6個字元"
              lazy-rules
              :rules="[
                val => 6 <= val.length || '密碼長度最小6個字元',
                val => val == this.rePassword || '密碼不一樣'
              ]"
              @input="
                formatVaild();
                passwordVaild();
              "
            />
          </div>
          <div class="column q-pa-sm">
            <q-input
              ref="rePasswordInput"
              outlined
              dense
              bg-color="grey-4"
              v-model="rePassword"
              type="password"
              label="確認密碼"
              :rules="[
                val => 6 <= val.length || '密碼長度最小6個字元',
                val => val == this.password || '密碼不一樣'
              ]"
              @input="
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
            @click="redirectSigninPage"
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
@import "../css/width.scss";

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

<script>
export default {
  name: "RegisterPage",
  data() {
    return {
      requireEmailAct: false,
      showName: "",
      loginName: "",
      email: "",
      emailRule: /^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z]+$/,
      password: "",
      rePassword: "",
      registerEnable: false,
      registerError: false,
      registerErrMsgArray: [],
      registerSuccess: false,
      registerSuccessMsg: "",
      id: -1,
      errList: []
    };
  },
  created() {
    this.getServerInfo();
  },
  methods: {
    passwordVaild() {
      if (this.passwordEqual() && this.password.length >= 6) {
        this.$refs.passwordInput.resetValidation();
        this.$refs.rePasswordInput.resetValidation();
      }
    },
    formatVaild() {
      if (1 > this.showName.length || this.showName.length > 15) {
        this.registerEnable = false;
        return false;
      }
      if (4 > this.loginName.length || this.loginName.length > 30) {
        this.registerEnable = false;
        return false;
      }

      if (this.email.length > 256 || this.email.search(this.emailRule) == -1) {
        this.registerEnable = false;
        return false;
      }
      if (!this.passwordEqual() || this.password.length < 6) {
        this.registerEnable = false;
        return false;
      }

      this.registerEnable = true;
      return true;
    },
    passwordEqual() {
      if (this.password == this.rePassword) {
        return true;
      } else {
        return false;
      }
    },
    redirectSigninPage() {
      this.$router.push("/signin");
    },
    redirectHomePage() {
      this.$router.push("/");
    },
    async getServerInfo() {
      var path = "/api/server-info";
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.requireEmailAct = data["require_email_activate"];
        })
        .catch(error => {
          if (error.request) {
          } else if (error.response) {
            switch (error.response.status) {
              //Internal Server Error
              case 500:
                break;
              //Unauthorized
              case 401:
                break;
            }
          }
        });
    },
    async sendRegisterData() {
      if (this.formatVaild()) {
        var path = "/api/register";
        await this.$axios
          .post(path, {
            show_name: this.showName,
            login_name: this.loginName,
            email: this.email,
            password: this.password
          })
          .then(res => {
            var data = res.data;
            this.id = data["id"];
            this.errList = data["err_list"];
            if (this.errList) {
              this.registerErrMsgArray = [];
              this.errList.forEach(errCode => {
                switch (errCode) {
                  case 1:
                    this.registerErrMsgArray.push("登入帳號已被使用過");
                    break;
                  case 2:
                    this.registerErrMsgArray.push("帳號長度錯誤");
                    break;
                  case 3:
                    this.registerErrMsgArray.push("顯示名稱長度錯誤");

                    break;
                  case 4:
                    this.registerErrMsgArray.push("密碼長度錯誤");
                    break;
                  case 5:
                    this.registerErrMsgArray.push(
                      "不支援大於256字元的電子郵件"
                    );

                    break;
                  case 6:
                    this.registerErrMsgArray.push("電子郵件格式錯誤");
                    break;
                  case 7:
                    this.registerErrMsgArray.push("電子郵件已被使用過");
                    break;
                  case 8:
                    this.registerErrMsgArray.push("系統錯誤");
                    break;
                  default:
                    break;
                }
              });
              this.registerError = true;
            } else {
              if (this.requireEmailAct) {
                this.registerSuccessMsg =
                  "認證信件已經寄到 " +
                  this.email +
                  " 請透過信中的連結完成註冊";
              }
              this.registerSuccess = true;
            }
          })
          .catch(error => {
            if (error) {
            }
          });
      } else {
      }
    }
  }
};
</script>
