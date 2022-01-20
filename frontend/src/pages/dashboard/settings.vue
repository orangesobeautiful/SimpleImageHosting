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

<script lang="ts">
import { ref, defineComponent } from 'vue';
import axios from 'axios';
import { api } from 'src/boot/axios';
import { json } from 'src/lib/common/type';

export default defineComponent({
  name: 'DashboardSettings',
  setup() {
    const showMsg = ref(false);
    const msgTitle = ref('');
    const msgContent = ref('');
    const hideSenderPwd = ref(true);
    //current value
    const hostname = ref('');
    const requireEmailAct = ref(false);
    const senderEmailServer = ref('');
    const senderEmailAddress = ref('');
    const senderEmailUser = ref('');
    // edit value
    const editHostname = ref('');
    const editRequireEmailAct = ref(false);
    const editSenderEmailServer = ref('');
    const editSenderEmailAddress = ref('');
    const editSenderEmailUser = ref('');
    const editSenderEmailPassword = ref('');

    // setting 的 json 格式
    interface settingJson {
      hostname: string;
      require_email_activate: boolean;
      sender_email_server: string;
      sender_email_address: string;
      sender_email_user: string;
      sender_email_password: string;
    }

    // getSettings 取得當前伺服器設定
    async function getSettings() {
      let path = 'dashboard/settings';
      await api.get(path).then((res) => {
        const data = res.data as settingJson;
        hostname.value = data['hostname'];
        requireEmailAct.value = data['require_email_activate'];
        senderEmailServer.value = data['sender_email_server'];
        senderEmailAddress.value = data['sender_email_address'];
        senderEmailUser.value = data['sender_email_user'];

        editHostname.value = hostname.value;
        editRequireEmailAct.value = requireEmailAct.value;
        editSenderEmailServer.value = senderEmailServer.value;
        editSenderEmailAddress.value = senderEmailAddress.value;
        editSenderEmailUser.value = senderEmailUser.value;

        if (senderEmailServer.value != '') {
        }
      });
    }
    void getSettings();

    async function editSettings() {
      let path = '/dashboard/settings';

      //要傳送的資料
      let patchData: json = {};
      patchData['hostname'] = editHostname.value;
      patchData['require_email_activate'] = editRequireEmailAct.value;
      if (
        senderEmailServer.value != editSenderEmailServer.value ||
        senderEmailAddress.value != editSenderEmailAddress.value ||
        senderEmailUser.value != editSenderEmailUser.value ||
        editSenderEmailPassword.value != editSenderEmailPassword.value
      ) {
        patchData['sender_email_server'] = editSenderEmailServer.value;
        patchData['sender_email_address'] = editSenderEmailAddress.value;
        patchData['sender_email_user'] = editSenderEmailUser.value;
        patchData['sender_email_password'] = editSenderEmailPassword.value;
      }

      await api
        .patch(path, patchData)
        .then((res) => {
          let data = res.data as json;
          msgTitle.value = '修改成功';
          msgContent.value = JSON.stringify(data, null, '\t');
          showMsg.value = true;
        })
        .catch((error) => {
          if (axios.isAxiosError(error)) {
            if (error.response) {
              msgTitle.value = error.response.status.toString();
              showMsg.value = true;
              msgContent.value = JSON.stringify(
                JSON.parse(error.response.data),
                null,
                '\t'
              );
            }
          }
        });
    }

    return {
      showMsg,
      msgTitle,
      msgContent,
      hideSenderPwd,
      //current value
      hostname,
      requireEmailAct,
      senderEmailServer,
      senderEmailAddress,
      senderEmailUser,
      // edit value
      editHostname,
      editRequireEmailAct,
      editSenderEmailServer,
      editSenderEmailAddress,
      editSenderEmailUser,
      editSenderEmailPassword,
      editSettings,
    };
  },
  methods: {},
});
</script>
