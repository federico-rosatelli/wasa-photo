<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			token_data: null,
			myUsername: ""
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.myUsername = "";
			try {
                //let username = myUsername
                let dataPost = {
                    username:this.myUsername
                }
				let response = await this.$axios.post("/signin",dataPost);
				this.token_data = response.data;
			} catch (e) {
				this.errormsg = "prova "+e.toString();
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
	<h1>{{this.token_data}}</h1>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
