<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			token: null,
            username: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.token = localStorage.getItem("Token")
				let response = await this.$axios.get("/",{headers:{"Token":this.token}});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
        async searchUser(){
            this.loading = true;
			this.errormsg = null;
			try {
				this.token = localStorage.getItem("Token")
				let response = await this.$axios.get("/search?query="+this.username,{headers:{"Token":this.token}});
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
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Search Profiles</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>
        <div class="login-form">
            <input type="text" placeholder="Username" v-model="username" >
			<input type="submit" value="Search" @click="searchUser">
        </div>
        <div v-for="item in this.some_data">
            <h1>{{item.Id}}</h1>
            <RouterLink v-bind:to="'/profile/'+item.Username" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#spike"/></svg>
								{{item.Username}}
							</RouterLink>
        </div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
