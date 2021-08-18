<template>
  <q-page>
    <q-dialog v-model="showMsg">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">{{ msgTitle }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none text-center">
          {{ msgContent }}
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="確認" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <div class="columns">
      <div class="column q-pa-sm">
        <q-input
          outlined
          dense
          bg-color="grey-4"
          v-model="editHostname"
          label="主機名稱"
          bottom-slots
          hint="網站所在的域名 (ex: www.exmple.com)"
          lazy-rules
        />
      </div>
      <q-checkbox
        v-model="editRequireEmailAct"
        label="註冊帳號需要email認證"
        color="cyan"
      />
      <div class="column q-pa-sm">
        <q-input
          outlined
          dense
          bg-color="grey-4"
          v-model="editSenderEmailServer"
          label="電子郵件伺服器"
          bottom-slots
          hint="需要包含 port (ex: mail.exmple.com:465)"
          lazy-rules
        />
      </div>
      <div class="column q-pa-sm">
        <q-input
          outlined
          dense
          bg-color="grey-4"
          v-model="editSenderEmailAddress"
          label="電子郵件位址"
          bottom-slots
          hint="email address"
          lazy-rules
        />
      </div>
      <div class="column q-pa-sm">
        <q-input
          outlined
          dense
          bg-color="grey-4"
          v-model="editSenderEmailUser"
          label="電子郵件使用者"
          bottom-slots
          hint="username"
          lazy-rules
        />
      </div>
      <div class="column q-pa-sm">
        <q-input
          outlined
          dense
          bg-color="grey-4"
          :type="hideSenderPwd ? 'password' : 'text'"
          v-model="editSenderEmailPassword"
          label="電子郵件密碼"
          bottom-slots
          hint="password"
          lazy-rules
          ><template v-slot:append>
            <q-icon
              :name="hideSenderPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="hideSenderPwd = !hideSenderPwd"
            /> </template
        ></q-input>
      </div>
      <div class="column q-pa-sm">
        <q-btn color="primary" label="修改" @click="editSettings" />
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped></style>

<script>
export default {
  name: "DashboardSettings",
  data() {
    return {
      showMsg: false,
      msgTitle: "",
      msgContent: "",
      hideSenderPwd: true,
      //current value
      hostname: "",
      requireEmailAct: false,
      senderEmailServer: "",
      senderEmailAddress: "",
      senderEmailUser: "",
      // edit value
      editHostname: "",
      editRequireEmailAct: false,
      editSenderEmailServer: "",
      editSenderEmailAddress: "",
      editSenderEmailUser: "",
      editSenderEmailPassword: ""
    };
  },
  async created() {
    await this.getSettings();
  },
  methods: {
    async getSettings() {
      var path = "/api/dashboard/settings";
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.hostname = data["hostname"];
          this.requireEmailAct = data["require_email_activate"];
          this.senderEmailServer = data["sender_email_server"];
          this.senderEmailAddress = data["sender_email_address"];
          this.senderEmailUser = data["sender_email_user"];

          this.editHostname = this.hostname;
          this.editRequireEmailAct = this.requireEmailAct;
          this.editSenderEmailServer = this.senderEmailServer;
          this.editSenderEmailAddress = this.senderEmailAddress;
          this.editSenderEmailUser = this.senderEmailUser;

          if (this.senderEmailServer != "") {
          }
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
                this.hasLogin = false;
                break;
            }
          }
        });
    },
    async editSettings() {
      var path = "/api/dashboard/settings";

      //要傳送的資料
      var patchData = {};
      patchData["hostname"] = this.editHostname;
      patchData["require_email_activate"] = this.editRequireEmailAct;
      if (
        this.senderEmailServer != this.editSenderEmailServer ||
        this.senderEmailAddress != this.editSenderEmailAddress ||
        this.senderEmailUser != this.editSenderEmailUser ||
        this.editSenderEmailPassword != this.editSenderEmailPassword
      ) {
        patchData["sender_email_server"] = this.editSenderEmailServer;
        patchData["sender_email_address"] = this.editSenderEmailAddress;
        patchData["sender_email_user"] = this.editSenderEmailUser;
        patchData["sender_email_password"] = this.editSenderEmailPassword;
      }

      await this.$axios
        .patch(path, patchData)
        .then(res => {
          var data = res.data;
          this.msgTitle = "修改成功";
          this.msgContent = JSON.stringify(data, null, "\t");
          this.showMsg = true;
        })
        .catch(error => {
          if (error.request) {
            this.msgTitle = error.request.status;
            this.msgContent = JSON.stringify(
              JSON.parse(error.request.response),
              null,
              "\t"
            );
            this.showMsg = true;
          } else if (error.response) {
          }
        });
    }
  }
};
</script>
