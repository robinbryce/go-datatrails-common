package redis

// There are 3 interfaces in this package:
//
//	1. hashstore - caching of entities
//	2. count limits - limit the number of entities allowed per tenant.
//	3. size limits - limit the size of an entity per tenant.
//
// 1. Hashstore
//
// Caches entities using a key/value store. Get, Set and Delete methods
// are available.
//
// 2. Count Limits
//
// Handles counting of resources and limiting the number of instances to a global
// value.
//
// The underlying counter is stored in REDIS using a serialised counter that
// is decremented whenever an instance of the entity is created. If the counter
// reaches zero, permission to create the entity is denied - Available() returns
// the number of available slots which in this case will be zero.
//
// Additionally the limit is fetched using a user-supplied method. Suitable logic to
// adjust the counter value when the limit changes is supplied in the initialise()
// method. Persistence is achieved by re-initialising when an error occurs accessing
// REDIS (in case the REDIS service is momentarily unavailable). Reiniialising consists
// of using a user supplied Counter() method to get the current number of entities.
//
// If the current counter value is unavailable from REDIS and we are unable to
// calculate the actual value using the Counter Method then we fail 'closed'.
//
// When creating an entity the caps limits is checked. There are three possibilities:
//
//    1. No limit is imposed i.e. 'unlimited' in which case the caps limit checking
//       is not used.
//    2. A limit is in place and the current number of entities is below the limit. In
//       which case the entity is created and - if no error occurs - the counter is
//       decremented.
//    3. A limit is in place and the current number of entities is at the limit. In
//       which case a 402 code is returned without attempting to create the entity.
//
// 3. Size Limits
//
// Limits the size of an entity - for example the maximum size of a blob.
//
