<template>
    <JsonModal id="eosModal" :title="$t('config.eos.title')" :description="$t('config.eos.description')"
        docs="/docs/reference/configuration/eos" endpoint="/config/eos" state-key="eos" data-testid="eos-modal"
        @changed="$emit('changed')">
        <template #default="{ values }">
            <FormRow id="eosServer" :label="$t('config.eos.labelServer')" example="localhost:8503">
                <input id="eosServer" v-model="values.url" class="form-control" required />
            </FormRow>

            <h6>{{ $t("config.eos.serverConfig") }}</h6>
            <FormRow id="yearlyConsumption" :label="$t('config.eos.yearlyConsumption')"
                :help="$t('config.eos.descriptionConsumption')" example="8000">
                <input id="yearlyConsumption" v-model="values.consumption" type="number" class="form-control" />
            </FormRow>
            <FormRow id="batteryMinSoc" :label="$t('config.eos.batteryMinSoc')" example="5">
                <input id="batteryMinSoc" v-model="values.batteryMinSoc" type="number" class="form-control" />
            </FormRow>
            <FormRow id="batteryMaxSoc" :label="$t('config.eos.batteryMaxSoc')" example="90">
                <input id="batteryMaxSoc" v-model="values.batteryMaxSoc" type="number" class="form-control" />
            </FormRow>
            <FormRow id="inverterMaxPower" :label="$t('config.eos.inverterMaxPower')" example="10000">
                <input id="inverterMaxPower" v-model="values.inverterMaxPower" type="number" class="form-control" />
            </FormRow>
            <FormRow id="batteryTariff" :label="$t('config.eos.batteryTariff')" example="0.08">
                <PropertyField id="batteryTariff" v-model="values.batteryTariff" type="Float" unit="€"
                    size="w-25 w-min-200" class="me-2" />
            </FormRow>
        </template>
    </JsonModal>
</template>

<script>
import JsonModal from "./JsonModal.vue";
import FormRow from "./FormRow.vue";
import PropertyField from "./PropertyField.vue";

export default {
    name: "EosModal",
    components: { FormRow, JsonModal, PropertyField },
    emits: ["changed"],
};
</script>