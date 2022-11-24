<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
                let username = myUsername
                let dataPost = {
                    username:username
                }
				let response = await this.$axios.post("/signin",dataPost);
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div class="login-form">
			<h1>Login Form</h1>
			    <input type="text" placeholder="Username" v-model="myUsername" >
				<input type="submit" value="Log In" @click="refresh">
			<div style="padding:15px">
		</div>
    </div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
