import com.owlike.genson.Genson;
import org.hyperledger.fabric.contract.Context;
import org.hyperledger.fabric.contract.ContractInterface;
import org.hyperledger.fabric.contract.annotation.*;
import org.hyperledger.fabric.shim.ChaincodeException;
import org.hyperledger.fabric.shim.ChaincodeStub;

import java.time.LocalDate;


@Contract(
        name = "Hash",
        info = @Info(
                title = "Hash contract for thesis",
                description = "A simple hash crud contract created for a thesis",
                version = "0.0.1-SNAPSHOT"))

@Default
public final class HashRepository implements ContractInterface{
    private final Genson genson = new Genson();

    @Transaction()
    public void initLedger(final Context ctx) {
        ChaincodeStub stub = ctx.getStub();
        Hash hash = new Hash("asdasdasdasd", LocalDate.now().toString(), "pippo");

        String agreementState = genson.serialize(hash);
        stub.putStringState("hash001", agreementState);
    }


    @Transaction()
    public Hash getHash(final Context ctx, final String key) {
        ChaincodeStub stub = ctx.getStub();
        String hashState = stub.getStringState(key);

        if (hashState.isEmpty()) {
            String errorMessage = String.format("Hash record %s does not exist", key);
            System.out.println(errorMessage);
            throw new ChaincodeException(errorMessage, "Hash not found");
        }

        Hash agreement = genson.deserialize(hashState, Hash.class);

        return agreement;
    }

    @Transaction()
    public Hash createHash(final Context ctx, final String key, final String hashString, final String date,
                                     final String user) {
        ChaincodeStub stub = ctx.getStub();

        String hashState = stub.getStringState(key);
        if (!hashState.isEmpty()) {
            String errorMessage = String.format("Hash record %s already exists", key);
            System.out.println(errorMessage);
            throw new ChaincodeException(errorMessage, "Hash already exists");
        }

        Hash hash = new Hash(hashString, date, user);
        hashState = genson.serialize(hash);
        stub.putStringState(key, hashState);

        return hash;
    }


    @Transaction()
    public Hash changeHashStatus(final Context ctx, final String key, final String date, final String user) {
        ChaincodeStub stub = ctx.getStub();

        String hashState = stub.getStringState(key);

        if (hashState.isEmpty()) {
            String errorMessage = String.format("Hash record %s does not exist", key);
            System.out.println(errorMessage);
            throw new ChaincodeException(errorMessage, "Hash not found");
        }

        Hash hash = genson.deserialize(hashState, Hash.class);

        Hash newHash = new Hash(hash.getHashString(),date,user);
        String newAgreementState = genson.serialize(newHash);
        stub.putStringState(key, newAgreementState);

        return newHash;
    }

}
