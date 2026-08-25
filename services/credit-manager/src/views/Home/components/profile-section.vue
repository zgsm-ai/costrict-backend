<template>
    <div class="profile-section">
        <div class="info-item">
            <div class="item-user-name">
                <div class="label">{{ t('homePage.userNameLabel') }}</div>
                <span :class="{ 'ml-1': !isZh }">{{ userName || '-' }}</span>
            </div>
            <div class="item-phone">
                <div class="label">{{ t('homePage.phoneLabel') }}</div>
                <span :class="{ 'ml-1': !isZh }">{{ phoneNumber || '-' }}</span>
            </div>
            <div
                class="item-account"
                v-if="!isPrivate"
            >
                <div class="label">
                    {{ t('homePage.githubAccount') }}
                </div>
                <span :class="{ 'ml-1': !isZh }">{{ githubName || '-' }}</span>
            </div>
            <div class="item-userId max-w-[100%]">
                <div class="label whitespace-nowrap">{{ t('homePage.userIdLabel') }}</div>
                <span
                    class="truncate flex-1"
                    :class="{ 'ml-1': !isZh }"
                    v-hover="userId"
                    >{{ userId || '-' }}</span
                >
                <n-icon
                    v-if="userId"
                    size="14"
                    style="margin-left: 8px; cursor: pointer"
                    color="#197DFF"
                    @click="handleCopyUserId"
                >
                    <copy-outline />
                </n-icon>
            </div>

            <div class="item-invite-code flex items-center">
                <div class="label opacity-70">
                    {{ t('homePage.inviteCodeLabel') }}
                </div>
                <span :class="{ 'ml-1': !isZh }">{{ inviteCode || '-' }}</span>
                <n-icon
                    v-if="inviteCode"
                    size="14"
                    style="margin-left: 8px; cursor: pointer"
                    color="#197DFF"
                    @click="handleCopyInviteCode"
                >
                    <copy-outline />
                </n-icon>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
/**
 * @file profile-section.vue
 * @description 个人信息组件 - 显示用户基本信息
 */
import { useI18n } from 'vue-i18n';
import { NIcon } from 'naive-ui';
import { CopyOutline } from '@vicons/ionicons5';

const { t } = useI18n();

// Props 定义
const props = withDefaults(
    defineProps<{
        userName?: string;
        phoneNumber?: string;
        githubName?: string;
        userId: string;
        inviteCode: string;
        isPrivate: boolean;
        isZh: boolean;
    }>(),
    {
        userName: undefined,
        phoneNumber: undefined,
        githubName: undefined,
        userId: '',
        inviteCode: '',
        isPrivate: false,
        isZh: false,
    },
);

// Emits 定义
const emit = defineEmits<{
    'copy-user-id': [value: string];
    'copy-invite-code': [value: string];
}>();

// 事件处理函数
const handleCopyUserId = () => {
    emit('copy-user-id', props.userId);
};

const handleCopyInviteCode = () => {
    emit('copy-invite-code', props.inviteCode);
};
</script>

<style scoped lang="less">
.profile-section {
    .info-item {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        column-gap: 40px;
        row-gap: 10px;

        .item-user-name,
        .item-phone {
            display: flex;
            align-items: center;

            .label {
                opacity: 0.7;
            }

            span {
                font-weight: 600;
            }
        }

        .item-account {
            display: flex;
            align-items: center;

            .label {
                opacity: 0.7;
            }

            span {
                font-weight: 600;
            }
        }

        .item-userId {
            display: flex;
            align-items: center;

            .label {
                opacity: 0.7;
            }
        }

        .item-invite-code {
            display: flex;
            align-items: center;

            .label {
                opacity: 0.7;
            }
        }
    }
}
</style>
