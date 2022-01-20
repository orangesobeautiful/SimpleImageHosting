<template>
  <q-page>
    <!-- loading page -->
    <loading-view :inputVisible="showLoading" />
    <!--normal page-->
    <div v-if="imageDataLoaded">
      <div class="row">
        <q-img fit="contain" class="image-block" :src="originalUrl" />
      </div>
      <div class="column q-pa-md q-gutter-md">
        <div class="row">
          <div class="row items-center">
            <q-btn class="q-mr-sm" round @click="push.userImagesPage(ownerID)">
              <q-avatar size="42px">
                <q-img :src="ownerAvatar" />
              </q-avatar>
            </q-btn>
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
    </div>

    <!--圖片修改對話框-->
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

    <!-- error page -->
    <div v-if="showErrPage" class="column err-colunm items-center">
      <div class="err-msg">
        <div class="text-center text-h5">
          {{ errTitle }}
        </div>
        <q-space style="height: 20px" />
        <div class="row justify-around err-redirect-link">
          <div class="text-h6 text-center" @click="push.homePage()">
            返回首頁
          </div>
          <div class="text-h6 text-center" @click="push.previousPage()">
            上一頁
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.image-block {
  max-height: calc(100vh - 50px);
  max-width: 100vw;
}
.owner-name-text {
  font-size: 15px;
  line-height: 40px;
}
.link-input {
  width: 368px;
}
.err-colunm {
  width: 100%;
}
.err-msg {
  margin-top: 24px;
  margin-bottom: 24px;
  width: 400px;
}
.err-redirect-link {
  color: #337f15;
}
</style>

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { copyToClipboard, useQuasar } from 'quasar';
import { api } from 'boot/axios';
import { useStore } from 'src/store';
import { Push } from 'src/lib/router/pushPage';
import { json } from 'src/lib/common/type';
import LoadingView from '../components/LoadingView.vue';

export default defineComponent({
  name: 'PageImage',
  components: {
    'loading-view': LoadingView,
  },
  setup() {
    const route = useRoute();
    const $q = useQuasar();
    const store = useStore();
    const push = new Push();

    const showLoading = ref(true);
    const showErrPage = ref(false);
    const imageDataLoaded = ref(false);
    const errTitle = ref('');
    const hashID = ref(route.params.imgHashID as string);
    const title = ref('');
    const ownerID = ref(-1);
    const ownerName = ref('');
    const ownerAvatar = ref('');
    const description = ref('');
    const originalUrl = ref('');
    const mdUrl = ref('');
    const host = ref(location.host);
    const createDateStr = ref('');
    interface imageLink {
      name: string;
      url: string;
    }
    const imageLinks = ref([] as imageLink[]);
    const tabRef = ref('about');
    const editTitle = ref('');
    const editDescription = ref('');
    const editDialog = ref(false);

    interface imageDataJson {
      title: string;
      owner_name: string;
      owner_id: number;
      owner_avatar: string;
      description: string;
      original_url: string;
      md_url?: string;
      create_at: number;
    }

    async function getImagesData() {
      let path = '/image/' + hashID.value;
      await api
        .get(path)
        .then((res) => {
          let data = res.data as imageDataJson;
          title.value = data['title'];
          ownerName.value = data['owner_name'];
          ownerID.value = data['owner_id'];
          ownerAvatar.value = data['owner_avatar'];
          description.value = data['description'];
          originalUrl.value = data['original_url'];
          if (data.md_url) {
            mdUrl.value = data['md_url'];
          }

          // 時間處理
          let create_millsec = data['create_at'] * 1000;
          let dateObject = new Date(create_millsec);
          let humanDateFormat = dateObject.toLocaleString();
          createDateStr.value = humanDateFormat;

          //連結處理
          let origin = location.origin;
          imageLinks.value.push({
            name: '圖片連結',
            url: origin + '/image/' + hashID.value,
          });
          imageLinks.value.push({
            name: '圖片URL',
            url: origin + '/' + originalUrl.value,
          });
          if (mdUrl.value != '') {
            imageLinks.value.push({
              name: '縮圖URL',
              url: origin + '/' + mdUrl.value,
            });
          }
          // 顯示正常頁面
          normalPage();
        })
        .catch((error) => {
          if (axios.isAxiosError(error)) {
            errTitle.value = '發生錯誤';
            if (error.response) {
              switch (error.response.status) {
                //Not Found
                case 404:
                  errTitle.value = '您訪問的連結錯誤或是圖片已被刪除';
                  break;
              }
            }
            // 顯示錯誤頁面
            errPage();
          }
        });
    }
    void getImagesData();

    async function editImageInfo() {
      var patchData: json = {};
      let datahasChanged = false;
      if (title.value != editTitle.value) {
        datahasChanged = true;
        patchData['title'] = editTitle.value;
      }
      if (description.value != editDescription.value) {
        datahasChanged = true;
        patchData['description'] = editDescription.value;
      }

      if (datahasChanged) {
        var path = '/image/' + hashID.value;

        await api
          .patch(path, patchData)
          .then(() => {
            $q.notify({
              message: '修改成功',
              color: 'green',
              timeout: 500,
              position: 'center',
              //persistent: true,
            });
            title.value = editTitle.value;
            description.value = editDescription.value;
          })
          .catch((error) => {
            console.log('has error');
            console.log(error);
            if (axios.isAxiosError(error)) {
              if (error.response) {
                console.log('error response');
                switch (error.response.status) {
                  //Internal Server Error
                  case 500:
                    $q.notify({
                      message: '伺服端錯誤',
                      color: 'red',
                      timeout: 1500,
                      position: 'center',
                    });
                    break;
                  //BadRequest
                  case 400:
                    $q.notify({
                      message: '錯誤的圖片ID',
                      color: 'red',
                      timeout: 1500,
                      position: 'center',
                    });
                    break;
                  //Unauthorized
                  case 401:
                    $q.notify({
                      message: '你需要登入',
                      color: 'red',
                      timeout: 1500,
                      position: 'center',
                    });
                    break;
                  //Forbidden
                  case 403:
                    $q.notify({
                      message: '權限不符',
                      color: 'red',
                      timeout: 1500,
                      position: 'center',
                    });
                    break;
                }
              }
            }
          });
      } else {
        $q.notify({
          message: '資料沒有變動',
          color: 'red',
          timeout: 1500,
          position: 'center',
        });
      }
    }

    async function deleteImage() {
      let path = '/image/' + hashID.value;
      await api
        .delete(path)
        .then(() => {
          $q.notify({
            message: '刪除成功',
            color: 'green',
            timeout: 500,
            position: 'center',
          });
          //回到個人介面
          push.userImagesPage(store.state.user.id);
        })
        .catch((error) => {
          if (axios.isAxiosError(error)) {
            if (error.response) {
              switch (error.response.status) {
                //Internal Server Error
                case 500:
                  $q.notify({
                    message: '伺服器內部錯誤',
                    color: 'red',
                    timeout: 1500,
                    position: 'center',
                  });
                  break;
                  break;
                //Unauthorized
                case 401:
                  $q.notify({
                    message: '你需要登入',
                    color: 'red',
                    timeout: 1500,
                    position: 'center',
                  });
                  break;
                //Forbidden
                case 403:
                  $q.notify({
                    message: '權限不符',
                    color: 'red',
                    timeout: 1500,
                    position: 'center',
                  });
                  break;
              }
            }
          }
        });
    }

    // 載入狀態處理 (loading, error ...)
    function errPage() {
      showLoading.value = false;
      showErrPage.value = true;
    }
    function normalPage() {
      showLoading.value = false;
      imageDataLoaded.value = true;
    }
    function showEditDialog() {
      editTitle.value = title.value;
      editDescription.value = description.value;
      editDialog.value = true;
    }

    function copyURL(linkIndex: number) {
      copyToClipboard(imageLinks.value[linkIndex].url)
        .then(() => {
          // success!
          $q.notify({
            message: '複製成功',
            color: 'green',
            timeout: 500,
          });
        })
        .catch(() => {
          // fail
          $q.notify({
            message: '複製失敗',
            color: 'red',
            timeout: 500,
          });
        });
    }

    function deleteConfirm() {
      $q.dialog({
        title: '刪除確認',
        message: '確定要刪除嗎?',
        ok: '刪除',
        cancel: '取消',
        focus: 'cancel',
      })
        .onOk(() => {
          void deleteImage();
        })
        .onCancel(() => {
          //do nothing
        });
    }

    return {
      push,
      showLoading,
      showErrPage,
      imageDataLoaded,
      errTitle,
      hashID,
      title,
      ownerID,
      ownerName,
      ownerAvatar,
      description,
      originalUrl,
      mdUrl,
      host,
      createDateStr,
      imageLinks,
      tabRef,
      editTitle,
      editDescription,
      editDialog,
      copyURL,
      showEditDialog,
      deleteConfirm,
      editImageInfo,
    };
  },
});
</script>
