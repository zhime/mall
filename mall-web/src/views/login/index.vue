<template>
  <div class="login-page">
    <div class="login-header">
      <h1>欢迎登录</h1>
      <p>Mall商城</p>
    </div>

    <div class="login-form">
      <van-cell-group inset>
        <van-field
          v-model="phone"
          type="tel"
          label="手机号"
          placeholder="请输入手机号"
          :rules="[{ validator: validatePhone, message: '请输入正确的手机号' }]"
        />
        <van-field
          v-model="code"
          type="digit"
          label="验证码"
          placeholder="请输入验证码"
        >
          <template #button>
            <van-button
              size="small"
              type="primary"
              :disabled="!canSendCode || sending"
              @click="sendCode"
            >
              {{ codeText }}
            </van-button>
          </template>
        </van-field>
      </van-cell-group>

      <div class="login-actions">
        <van-button
          type="primary"
          size="large"
          block
          :loading="logging"
          :disabled="!canLogin"
          @click="handleLogin"
        >
          登录
        </van-button>
      </div>

      <div class="login-tips">
        <p>登录即表示同意<span class="link">《用户协议》</span>和<span class="link">《隐私政策》</span></p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { isPhone } from '@/utils'
import { useUserStore } from '@/store/user'
import { sendSmsCode, loginByPhone } from '@/api/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const phone = ref('')
const code = ref('')
const sending = ref(false)
const logging = ref(false)
const countdown = ref(0)
let timer: NodeJS.Timeout | null = null

const canSendCode = computed(() => isPhone(phone.value))
const canLogin = computed(() => canSendCode.value && code.value.length === 6)

const codeText = computed(() => {
  if (countdown.value > 0) {
    return `${countdown.value}s后重发`
  }
  return sending.value ? '发送中...' : '获取验证码'
})

const validatePhone = (value: string) => isPhone(value)

// 发送验证码
const sendCode = async () => {
  if (!canSendCode.value || sending.value) return
  
  try {
    sending.value = true
    await sendSmsCode(phone.value, 'login')
    showToast('验证码发送成功')
    
    // 开始倒计时
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer!)
        timer = null
      }
    }, 1000)
  } catch (error: any) {
    showToast(error.message || '发送失败')
  } finally {
    sending.value = false
  }
}

// 登录
const handleLogin = async () => {
  if (!canLogin.value || logging.value) return
  
  try {
    logging.value = true
    const { data } = await loginByPhone({
      phone: phone.value,
      code: code.value,
      type: 'sms'
    })
    
    userStore.login(data.token, data.user)
    showToast('登录成功')
    
    // 重定向到目标页面
    const redirect = (route.query.redirect as string) || '/'
    router.replace(redirect)
  } catch (error: any) {
    showToast(error.message || '登录失败')
  } finally {
    logging.value = false
  }
}

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #ee0a24 0%, #ff6034 100%);
  padding: 60px 20px 20px;

  .login-header {
    text-align: center;
    color: white;
    margin-bottom: 60px;

    h1 {
      font-size: 28px;
      font-weight: 300;
      margin: 0 0 8px 0;
    }

    p {
      font-size: 16px;
      opacity: 0.8;
      margin: 0;
    }
  }

  .login-form {
    .login-actions {
      padding: 24px 16px;
    }

    .login-tips {
      padding: 0 16px;
      text-align: center;

      p {
        font-size: 12px;
        color: #666;
        line-height: 1.5;

        .link {
          color: #ee0a24;
        }
      }
    }
  }
}
</style>