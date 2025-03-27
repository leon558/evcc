<template>
	<JsonModal
		id="eosModal"
		:title="$t('config.eos.title')"
		:description="$t('config.eos.description')"
		docs="/docs/reference/configuration/eos"
		endpoint="/config/eos"
		state-key="eos"
		data-testid="eos-modal"
		@changed="$emit('changed')"
	>
		<template #default="{ values }">
			<FormRow
				id="eosServer"
				:label="$t('config.eos.labelServer')"
				example="localhost:8503"
			>
				<input id="eosServer" v-model="values.url" class="form-control" required />
			</FormRow>

			<h6>{{ $t("config.eos.serverConfig") }}</h6>
			<FormRow
				id="yearlyConsumption"
				:label="$t('config.eos.yearlyConsumption')"
				:help="$t('config.eos.descriptionConsumption')"
				example="8000"
			>
				<input id="yearlyConsumption" type="number" v-model="values.consumption" class="form-control" />
			</FormRow>
		</template>
	</JsonModal>
</template>

<script>
import JsonModal from "./JsonModal.vue";
import FormRow from "./FormRow.vue";

export default {
	name: "EosModal",
	components: { FormRow, JsonModal },
	emits: ["changed"],
};
</script>