<template>
    <div class="v-select">
        <div class="down-icon" :style="display?'rotate: 180deg;':''">
            <Icon name="ion:chevron-down-outline" color="#989898" size="20" />
        </div>
        <div class="value" @click="display = !display">
            <div class="item-selected m-2">
                <span>Vue</span>
                <span>
                    <Icon name="ion:close" color="#989898" size="20" />
                </span>
            </div>
            <div class="item-selected m-2">
                <span class="me-2">后端</span>
                <span class="v-icon">
                    <Icon name="ion:close" color="#989898" size="20" />
                </span>
            </div>
        </div>
        <div class="options" :class="display? 'd-inline-flex':'d-none'">
            <div class="form-check m-2">
                <input class="form-check-input" type="checkbox" value="" id="options-1">
                <label class="form-check-label" for="options-1">
                    Go
                </label>
            </div>

            <div class="form-check m-2">
                <input class="form-check-input" type="checkbox" value="2" id="options-2">
                <label class="form-check-label" for="options-2">
                    Vue
                </label>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const display = ref<boolean>(false)

</script>

<style scoped lang="scss">
.v-select {
    position: reactive;
    border: 1px solid #dee2e6;
    border-radius: 4px;
    cursor: pointer;

    .down-icon {
        position: absolute;
        right: 10px;
        top: 50%;
        margin-top: -12px;
    }

    .item-selected {
        display: inline-flex;
        background: #e2eae2;
        font-size: 14px;
        padding: 4px 8px;
        border-radius: 4px;
    }

    .v-icon {
        cursor: pointer;
    }

    .options {
        // display: none; // d-inline-flex
        transition: all 0.5;
        position: absolute;
        bottom: calc(-100% - 5px);
        z-index: 9;
        width: 100%;
        border: 1px solid #e2eae2;
        background: #fff;
        border-radius: 4px;
        padding: 0 20px;
        color: #666;
    }
}</style>