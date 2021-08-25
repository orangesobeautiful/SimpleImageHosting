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

        <q-card-actions align="center">
          <q-btn flat label="確認" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <div class="row">
      <q-img
        contain
        class="image-block"
        :src="originalUrl"
        native-context-menu
      />
    </div>
    <div class="column q-pa-md q-gutter-md">
      <div class="row">
        <div class="row items-center">
          <q-avatar class="q-mr-sm" size="40px">
            <q-img :src="ownerAvatar" />
          </q-avatar>
          <div class="row">
            <div class="owner-name-text q-mr-md">{{ ownerName }}</div>
            <div
              v-if="this.$store.state.user.id == this.ownerID"
              class="row items-center"
            >
              <q-btn
                flat
                text-color="blue"
                class="q-mr-sm"
                @click="showEditDialog"
                ><q-icon name="edit" />修改</q-btn
              >
              <q-btn
                flat
                text-color="red"
                class="q-mr-sm"
                @click="deleteConfirm"
                ><q-icon name="delete" />刪除</q-btn
              >
            </div>
          </div>
        </div>
      </div>
      <div class="text-h6 text-weight-bold">{{ title }}</div>
      <q-tabs v-model="tabRef" align="left">
        <q-tab name="about" label="資訊" />
        <q-tab name="link" label="連結" />
      </q-tabs>
      <q-tab-panels v-model="tabRef" animated class="q-gutter-none">
        <q-tab-panel name="about">
          <div class="column q-gutter-md">
            <div class="row">
              <div>上傳時間：</div>
              <div>{{ createDateStr }}</div>
            </div>
            <div class="text-body1">{{ description }}</div>
          </div>
        </q-tab-panel>
        <q-tab-panel name="link">
          <div class="column">
            <div v-for="(item, index) in imageLinks" :key="index">
              <div class="row items-center q-mb-md">
                <div class="q-mr-sm text-subtitle1">
                  {{ item.name }}
                </div>

                <q-input
                  v-model="item.url"
                  bg-color="grey-3"
                  class="link-input"
                  filled
                  dense
                >
                  <template v-slot:append>
                    <q-btn
                      dense
                      size="md"
                      color="grey-8"
                      @click="copyURL(index)"
                      >&emsp;複製&emsp;</q-btn
                    >
                  </template>
                </q-input>
              </div>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </div>
    <q-dialog v-model="editDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">編輯圖片</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="text-subtitle1">標題</div>
          <q-input dense v-model="editTitle" autofocus filled clearable />
          <div class="text-subtitle1">內容描述</div>
          <q-input v-model="editDescription" filled clearable type="textarea" />
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="取消" v-close-popup />
          <q-btn flat label="修改" v-close-popup @click="editImageInfo" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<style lang="sass" scoped>
.image-block
  max-height: calc(100vh - 50px)
  max-width: 100vw
.owner-name-text
    font-size: 15px
    line-height: 40px
.link-input
    width: 368px
</style>

<script>
import { copyToClipboard } from "quasar";

export default {
  name: "PageImage",
  data() {
    return {
      showMsg: false,
      msgTitle: "",
      msgContent: "",
      hashID: this.$route.params.imgHashID,
      title: "",
      ownerID: "",
      ownerName: "",
      ownerAvatar: "",
      description: "",
      originalUrl: "",
      mdUrl: "",
      host: location.host,
      createDateStr: "",
      imageLinks: Array(),
      tabRef: "about",
      editTitle: "",
      editDescription: "",
      editDialog: false
    };
  },
  created() {
    this.getImagesData();
  },
  methods: {
    copyURL(linkIndex) {
      copyToClipboard(this.imageLinks[linkIndex].url)
        .then(() => {
          // success!
        })
        .catch(() => {
          // fail
        });
      this.$q.notify({
        message: "複製成功",
        color: "green",
        timeout: "500"
      });
    },
    async getImagesData() {
      var path = "/api/image/" + this.hashID;
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.title = data["title"];
          this.ownerName = data["owner_name"];
          this.ownerID = data["owner_id"];
          this.ownerAvatar = data["owner_avatar"];
          this.description = data["description"];
          this.originalUrl = data["original_url"];
          if (data.md_url != null) {
            this.mdUrl = data["md_url"];
          }

          // 時間處理
          var create_millsec = data["create_at"] * 1000;
          var dateObject = new Date(create_millsec);
          var humanDateFormat = dateObject.toLocaleString();
          this.createDateStr = humanDateFormat;

          //連結處理
          var origin = location.origin;
          this.imageLinks.push({
            name: "圖片連結",
            url: origin + "/image/" + this.hashID
          });
          this.imageLinks.push({
            name: "圖片URL",
            url: origin + "/" + this.originalUrl
          });
          if (this.mdUrl != "") {
            this.imageLinks.push({
              name: "縮圖URL",
              url: origin + "/" + this.mdUrl
            });
          }
        })
        .catch(error => {
          if (error.request) {
            switch (error.response.status) {
              //Internal Server Error
              case 500:
                break;
              //Unauthorized
              case 401:
                break;
              //Not Found
              case 404:
                this.msgTitle = "圖片不存在";
                this.msgContent = "您訪問的連結錯誤或是圖片已被刪除";
                this.showMsg = true;
                break;
            }
          } else if (error.response) {
          }
        });
    },
    async deleteConfirm() {
      this.$q
        .dialog({
          title: "刪除確認",
          message: "確定要刪除嗎?",
          ok: "刪除",
          cancel: "取消",
          focus: "cancel"
        })
        .onOk(() => {
          this.deleteImage();
        })
        .onCancel(() => {
          //do nothing
        });
    },
    async deleteImage() {
      var path = "/api/image/" + this.hashID;
      await this.$axios
        .delete(path)
        .then(() => {
          this.$q.notify({
            message: "刪除成功",
            color: "green",
            timeout: "500",
            position: "center"
          });
          //回到個人介面
          this.$router.push("/user/" + this.$store.state.user.id + "/images");
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
                this.$q.notify({
                  message: "你需要登入",
                  color: "red",
                  timeout: "1500",
                  position: "center"
                });
                break;
              //Forbidden
              case 403:
                this.$q.notify({
                  message: "權限不符",
                  color: "red",
                  timeout: "1500",
                  position: "center"
                });
                break;
            }
          }
        });
    },
    showEditDialog() {
      this.editTitle = this.title;
      this.editDescription = this.description;
      this.editDialog = true;
    },
    async editImageInfo() {
      var putData = {};
      var datahasChanged = false;
      if (this.title != this.editTitle) {
        datahasChanged = true;
        putData["title"] = this.editTitle;
      }
      if (this.description != this.editDescription) {
        datahasChanged = true;
        putData["description"] = this.editDescription;
      }

      if (datahasChanged) {
        var path = "/api/image/" + this.hashID;

        await this.$axios
          .patch(path, putData)
          .then(() => {
            this.$q.notify({
              message: "修改成功",
              color: "green",
              persistent: true,
              timeout: "500",
              position: "center"
            });
            this.title = this.editTitle;
            this.description = this.editDescription;
          })
          .catch(error => {
            console.log("has error");
            console.log(error);
            if (error.request) {
              console.log("error request");
            } else if (error.response) {
              console.log("error response");
              switch (error.response.status) {
                //Internal Server Error
                case 500:
                  this.$q.notify({
                    message: "伺服端錯誤",
                    color: "red",
                    timeout: "1500",
                    position: "center"
                  });
                  break;
                //BadRequest
                case 400:
                  this.$q.notify({
                    message: "錯誤的圖片ID",
                    color: "red",
                    timeout: "1500",
                    position: "center"
                  });
                  break;
                //Unauthorized
                case 401:
                  this.$q.notify({
                    message: "你需要登入",
                    color: "red",
                    timeout: "1500",
                    position: "center"
                  });
                  break;
                //Forbidden
                case 403:
                  this.$q.notify({
                    message: "權限不符",
                    color: "red",
                    timeout: "1500",
                    position: "center"
                  });
                  break;
              }
            }
          });
      } else {
        this.$q.notify({
          message: "資料沒有變動",
          color: "red",
          timeout: "1500",
          position: "center"
        });
      }
    }
  }
};
</script>
