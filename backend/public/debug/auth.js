var uiConfig = {
  signInSuccessUrl: '/debug/',
  signInOptions: [
    {
      provider: firebase.auth.GoogleAuthProvider.PROVIDER_ID,
      scopes: [
        'profile'
      ],
      customParameters: {
        prompt: 'select_account'
      },
    },
  ],
  signInFlow: 'popup',
  callbacks: {
    signInSuccessWithAuthResult: function(authResult, redirectUrl) {
      return true;
    },
  },
};
firebase.auth().onAuthStateChanged(function(user) {
  if (user) {
    document.getElementById('firebaseui-auth-container').innerHTML = user.displayName + "<br>" + user.email + "<br>" + "<button onclick='signOut()'>Sign Out</button><br>";
  } else {
    var ui = new firebaseui.auth.AuthUI(firebase.auth());
    ui.start('#firebaseui-auth-container', uiConfig);
  }
});
function signOut() {
  firebase.auth().signOut().then(function() {
    location.href = '/debug/auth.html';
  });
}
