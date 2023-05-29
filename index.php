if (@$_REQUEST['google-login']){

		$data = json_decode($_REQUEST['google-login']);
		if (!$data){
			$data = json_decode( stripslashes( $_REQUEST['google-login']));
		}

		if (!@$data->idToken || !@$data->user){
			echo json_encode(['error'=>'missing required params']);
			exit;
		}

		$verifyUrl = "https://www.googleapis.com/oauth2/v1/tokeninfo?id_token=$data->idToken";
		$a = ArzUtility::curll($verifyUrl);

		if (@$a['error']){
			echo json_encode(['error'=>'cant get response from google']);
			exit;
		}

		$b = json_decode($a['content']);
		if (!$b){
			echo json_encode(['error'=>'cant parse response from google']);
			exit;
		}
		else if (@$b->error){
			echo json_encode(['error'=>"google says: $b->error"]);
			exit;
		}

		if (@$b->email != @$data->user->email){
			echo json_encode(['error'=>"invalid google token"]);
			exit;
		}

		require_once __DIR__ .'/../../plugins/arz-profile/templates/rewrite/signFuncs.php';

		global $wpdb;
		$email = $data->user->email;
		if (get_user_by('email' , $email)) {
            $user = get_user_by('email' , $email);
            $user_agent = arzdigital_google_login_get_browser();
            $insert_session_for_google_login = $wpdb->insert('wp_session', array(
                'user_id' => $user->ID,
                'status' => 'login',
                'token' => 'google_login_user',
                'browser' => $user_agent['browser'],
                'os' => $user_agent['platform'],
                'time' => time()
            ));

            $userid = $user->ID;

            if ($insert_session_for_google_login) {
                arzdigital_google_login_login_the_user($user->user_login, $user->ID,0);
            }

			logUserLogins(['uid'=>$userid,'from'=>'app-google-login']);
        } else {

            $userdata = array(
                'user_login' => generateRandomString(),
                'display_name' => @$data->user->name,
                'user_email' => $email,
                'user_pass' => arzdigital_google_login_generate_password(),
            );
            $insert_in_wp_users = wp_insert_user($userdata);
            $insert_in_google_table = $wpdb->insert('wp_users_logged_in_by_google', array(
                'google_id' => @$data->user->id,
                'user_id' => $insert_in_wp_users,
                'name' => @$data->user->name,
                'email' => $email,
                'profile_image' => @$data->user->photo
            ));

            $userid = $insert_in_wp_users;
            update_user_meta($insert_in_wp_users, 'arz_mobile', 0);
            update_user_meta($insert_in_wp_users, 'user_phone_verify', 0);
            update_user_meta($insert_in_wp_users, 'user_email_verify', time());
			update_user_meta($insert_in_wp_users, 'arz_registerApp', 'app');
            $user_login = get_user_by('id', $insert_in_wp_users);
            if ($insert_in_wp_users && $insert_in_google_table) {
                arzdigital_google_login_login_the_user($user_login->user_login, $insert_in_wp_users,0);

				logUserLogins(['uid'=>$userid,'from'=>'app-google-register']);
            }
			else {
				echo json_encode(['error'=>"Sign up failed"]);
				exit;
			}



        }




		$token = generateJwtToken($userid);
		echo json_encode(['ok'=>$token]);



        exit;
	}
	else if (@$_REQUEST['apple-login']){

		$data = json_decode($_REQUEST['apple-login']);
		if (!$data){

			$data = json_decode( stripslashes( $_REQUEST['apple-login']));
		}

		if (!@$data->identityToken ){
			echo json_encode(['error'=>'missing required params']);
			exit;
		}


		require_once APPLE_LIBPATH .'/vendor/autoload.php';

		$accessToken = '';
		try {
			$params = json_decode(APPLE_PARAMS);
			Firebase\JWT\JWT::$leeway = 60;
			$provider = new League\OAuth2\Client\Provider\Apple([
				'clientId'          => $params->clientId,
				'teamId'            => $params->teamId,
				'keyFileId'         => $params->keyFileId,
				'keyFilePath'       => APPLE_LIBPATH . "/AuthKey_{$params->keyFileId}.p8",
				'redirectUri'       => $params->redirectUri,
			]);

			$token = $provider->getAccessToken('authorization_code', [
				'code' => $data->authorizationCode
			]);

			$accessToken = $token->getToken();
			$verified = true;

		}
		catch(Exception $e){
			$verified = false;
		}

		if (!$verified){
			echo json_encode(['error'=>'invalid apple jwt']);
			exit;
		}

		$a = explode('.',$data->identityToken);
		$jwtParsed = json_decode(base64_decode($a[1]));

		$email = @$jwtParsed->email;
		//-----------------------------------


		$appleUid = $jwtParsed->sub;
		global $wpdb;
		$exists = $wpdb->get_row($wpdb->prepare("select * from wp_usermeta where meta_key='apple_uid' and meta_value=%s limit 1",[ $appleUid ]));


		$userid = 0;
		$newUser = false;
		if ($exists){

			$user = get_user_by('id' , $exists->user_id);

			if ($user){
				$userid = $user->ID;
			}
			else $newUser = true;

		}
		else {

			$user = get_user_by('email' , $email);
			if ($user){
				$userid = $user->ID;

				logUserLogins(['uid'=>$userid,'from'=>'app-apple-login']);
			}
			else $newUser = true;

		}

		if ($newUser) {

			if (!@$_REQUEST['name'] ){
				echo json_encode(['error'=>'name missing']);
				exit;
			}

			require_once __DIR__ .'/../../plugins/arz-profile/templates/rewrite/signFuncs.php';

			$name =  $_REQUEST['name'];
            $userdata = array(
                'user_login' => generateRandomString(),
                'display_name' => $name,
                'user_email' => $email,
                'user_pass' => arzdigital_google_login_generate_password(),
            );
            $insert_in_wp_users = wp_insert_user($userdata);
			/*
            $insert_in_google_table = $wpdb->insert('wp_users_logged_in_by_google', array(
                'google_id' => @$data->user->id,
                'user_id' => $insert_in_wp_users,
                'name' => @$data->user->name,
                'email' => $email,
                'profile_image' => @$data->user->photo
            ));
			*/
            $userid = $insert_in_wp_users;
            update_user_meta($insert_in_wp_users, 'arz_mobile', 0);
            update_user_meta($insert_in_wp_users, 'user_phone_verify', 0);
            update_user_meta($insert_in_wp_users, 'user_email_verify', time());
			update_user_meta($insert_in_wp_users, 'arz_registerApp', 'app');

			update_user_meta($insert_in_wp_users, 'apple_uid', $appleUid);
            /*
			$user_login = get_user_by('id', $insert_in_wp_users);

            if ($insert_in_wp_users && $insert_in_google_table) {
                arzdigital_google_login_login_the_user($user_login->user_login, $insert_in_wp_users,0);
            }
			else {
				echo json_encode(['error'=>"Sign up failed"]);
				exit;
			}
			*/

			logUserLogins(['uid'=>$userid,'from'=>'app-apple-register']);
        }

		//update_user_meta($userid,'apple-access-token',$accessToken);

		$token = generateJwtToken($userid);
		$token['apple_access_token'] = $accessToken;
		echo json_encode(['ok'=>$token]);



        exit;
	}