#ifndef RUST_C_IDENTITY_MANAGER
#define RUST_C_IDENTITY_MANAGER

#include <stdint.h>

typedef struct IdentityManager identity_manager_t;
typedef struct CCredential {
    char *cred;
}credential_t;

//EXPORTED FUNCTIONS
extern credential_t const *new_ccredential(char const *cred);
extern void drop_ccredential(credential_t *);

extern identity_manager_t *new_identity_manager(int mainnet, int persistence, char const *folder_dir, char const *vault_psw);
extern void drop_identity_manager(identity_manager_t *);
extern char const *create_identity(identity_manager_t *, char const *);
extern char const *get_identity_did(identity_manager_t *, char const *);
extern int store_credential(identity_manager_t *, char const *, credential_t const *);
extern credential_t *get_credential(identity_manager_t *, char const *);
extern void drop_str(char *);


#endif //RUST_C_IDENTITY_MANAGER
