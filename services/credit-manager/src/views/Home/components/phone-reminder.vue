<template>
    <transition name="phone-reminder-fade">
        <div v-if="visible" class="phone-reminder">
            <div class="reminder-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="#2080F0" stroke-width="1.5" />
                    <line x1="12" y1="10" x2="12" y2="17" stroke="#2080F0" stroke-width="2" stroke-linecap="round" />
                    <circle cx="12" cy="7" r="1" fill="#2080F0" />
                </svg>
            </div>
            <span class="reminder-text">
                {{ text }}
            </span>
            <span v-if="action" class="reminder-action" @click="handleBind">
                {{ action }} →
            </span>
            <button class="reminder-close" @click="handleClose" :aria-label="closeLabel">
                <svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                >
                    <line x1="18" y1="6" x2="6" y2="18" />
                    <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
            </button>
        </div>
    </transition>
</template>

<script setup lang="ts">
/**
 * @file phone-reminder.vue
 * @description 未绑定手机号提醒条幅 - 自定义实现，匹配项目暗黑主题风格
 */
import { ref } from 'vue';

defineProps<{
    text: string;
    action?: string;
    closeLabel?: string;
}>();

const emit = defineEmits<{
    bind: [];
    close: [];
}>();

const visible = ref(true);

const handleBind = () => {
    emit('bind');
};

const handleClose = () => {
    visible.value = false;
    emit('close');
};
</script>

<style scoped lang="less">
.phone-reminder {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border-radius: 10px;
    background: linear-gradient(
        135deg,
        rgba(24, 118, 242, 0.08) 0%,
        rgba(24, 118, 242, 0.03) 50%,
        rgba(255, 255, 255, 0.02) 100%
    );
    border: 1px solid rgba(24, 118, 242, 0.15);
    border-left: 3px solid #1876f2;
    backdrop-filter: blur(10px);
    box-shadow:
        0 2px 12px rgba(0, 0, 0, 0.08),
        inset 0 1px 0 rgba(255, 255, 255, 0.03);
    transition: opacity 0.3s ease, transform 0.3s ease;
    position: relative;
    overflow: hidden;

    &::before {
        content: '';
        position: absolute;
        top: -50%;
        left: -20%;
        width: 60%;
        height: 200%;
        background: radial-gradient(
            ellipse at center,
            rgba(24, 118, 242, 0.04) 0%,
            transparent 70%
        );
        pointer-events: none;
    }

    .reminder-icon {
        flex-shrink: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 20px;
        height: 20px;
    }

    .reminder-text {
        color: rgba(255, 255, 255, 0.8);
        font-size: 13px;
        line-height: 1.5;
    }

    .reminder-action {
        flex-shrink: 0;
        font-size: 13px;
        color: #1876f2;
        cursor: pointer;
        white-space: nowrap;
        transition: opacity 0.2s ease;

        &:hover {
            opacity: 0.8;
        }
    }

    .reminder-close {
        flex-shrink: 0;
        margin-left: auto;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
        border: none;
        background: transparent;
        color: rgba(255, 255, 255, 0.3);
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s ease;
        padding: 0;

        &:hover {
            color: rgba(255, 255, 255, 0.8);
            background: rgba(255, 255, 255, 0.08);
        }

        &:active {
            background: rgba(255, 255, 255, 0.12);
        }
    }
}

.phone-reminder-fade-leave-active {
    transition: opacity 0.25s ease, transform 0.25s ease;
}

.phone-reminder-fade-leave-to {
    opacity: 0;
    transform: translateY(-4px);
}
</style>
