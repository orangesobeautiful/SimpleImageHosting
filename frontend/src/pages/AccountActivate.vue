<template>
  <q-page class="flex flex-center">
    <q-dialog v-model="showMsg">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">{{ msgTitle }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none text-center">
          {{ msgContent }}
        </q-card-section>

        <q-card-actions align="center">
          <q-btn
            flat
            label="確認"
            color="primary"
            v-close-popup
            @click="confirmBtnClick"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
<script>
export default {
  name: "AccountActivatePage",
  data() {
    return {
      showMsg: false,
      actSuccess: false,
      msgTitle: "",
      msgContent: "",
      userID: 0,
      grade: 0,
      showName: ""
    };
  },
  created() {
    this.activateAccount();
  },
  methods: {
    async activateAccount() {
      var path = "/api/account-activate/" + this.$route.params.token;
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.userID = data["user_id"];
          this.grade = data["grade"];
          this.showName = data["show_name"];
          this.actSuccess = true;
          this.msgTitle = "認證成功";
          this.msgContent = "歡迎您的加入 " + this.showName;
          this.showMsg = true;
        })
        .catch(error => {
          if (error.request) {
            switch (error.request.status) {
              //Internal Server Error
              case 500:
                this.msgTitle = "認證失敗";
                this.msgContent = "伺服器內部錯誤";
                this.showMsg = true;
                break;
              //Forbidden
              case 403:
                this.msgTitle = "認證失敗";
                this.msgContent = "錯誤或已失效的連結";
                this.showMsg = true;
                break;
            }
          } else if (error.response) {
            this.msgTitle = "認證失敗";
            this.msgContent = error.response.data;
            this.showMsg = true;
          }
        });
    },
    confirmBtnClick() {
      if (this.actSuccess) {
        this.$router.push("/user/" + this.userID.toString() + "/images");
      } else {
        this.$router.push("/");
      }
    }
  }
};
</script>
